package server

import (
	"context"
	"crypto/rand"
	"encoding/json"
	stderrors "errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/psyhatter/word-of-wisdom/pkg/logger"
	"github.com/psyhatter/word-of-wisdom/pkg/models"
	"github.com/psyhatter/word-of-wisdom/pkg/puzzle"
	"go.uber.org/zap"
)

// DDOSProtector represents a struct used for handling DDOS protection.
type DDOSProtector struct {
	// Timeout represents the duration within which clients must send proof-of-work.
	Timeout time.Duration

	// NonceSize represents the size of the nonce in bytes.
	NonceSize int

	// Complexity represents the complexity level of a task. It is an integer value
	// between 0 and 255.
	Complexity uint8
}

// Handle handles the DDOS protection logic, and returns an error if any. It
// encodes the puzzle into JSON and sends it to the client. Then, it decodes the
// solution from JSON and checks if it meets the required difficulty level.
// Finally, it returns the result of the solution check.
func (d DDOSProtector) Handle(ctx context.Context, conn net.Conn) (err error) {
	enc := json.NewEncoder(conn)

	defer func() {
		resp := models.VerificationResponse{Success: err == nil}
		if err != nil {
			resp.ErrorMessage = err.Error()
		}

		err = JoinErrors(err, enc.Encode(resp))
	}()

	p := models.Puzzle{
		Difficulty: int(d.Complexity),
		Nonce:      make([]byte, d.NonceSize),
	}

	if d.Timeout > 0 {
		p.Deadline = time.Now().Add(d.Timeout)

		if err = conn.SetDeadline(p.Deadline); err != nil {
			return errors.Wrap(err, "setting deadline")
		}

		defer func() { err = JoinErrors(err, errors.WithStack(conn.SetDeadline(time.Time{}))) }()
	}

	if _, err = rand.Read(p.Nonce); err != nil {
		return errors.Wrap(err, "creating nonce")
	}

	if err = enc.Encode(p); err != nil {
		return errors.Wrap(err, "encoding puzzle")
	}

	var solution models.Solution
	err = json.NewDecoder(conn).Decode(&solution)
	if err != nil {
		return errors.Wrap(err, "reading solution")
	}

	logger.FromCtx(ctx).Debugw("solution", zap.Any("puzzle", p), zap.ByteString("solution", solution))

	return errors.WithStack(puzzle.CheckSolution(p, solution))
}

// Server represents a struct used for handling server configuration and logic.
type Server struct {
	// Address represents the server address for listening.
	Address string

	// DDOSProtector represents a struct used for handling DDOS protection.
	DDOSProtector DDOSProtector

	l net.Listener

	isClosed atomic.Bool
	wg       sync.WaitGroup
}

// Listen listens for incoming connections on the specified address and handles them using the provided handler function.
// It starts a goroutine that continuously accepts new connections and launches goroutines to handle each connection.
// If an error occurs while accepting a connection, it logs an error message and continues accepting new connections.
// After starting the goroutine, it returns nil immediately.
// If the context is cancelled, the goroutine will exit and all connections will be closed.
func (s *Server) Listen(ctx context.Context, handler func(ctx context.Context, conn net.Conn)) (err error) {
	s.l, err = (&net.ListenConfig{}).Listen(ctx, "tcp", s.Address)
	if err != nil {
		return errors.WithStack(err)
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		for {
			conn, err := s.l.Accept()
			if err != nil {
				if s.isClosed.Load() {
					return
				}

				logger.FromCtx(ctx).Errorw("accepting connection", zap.Error(err))
				continue
			}

			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				defer conn.Close()

				remoteAddr := conn.RemoteAddr()

				ctx := logger.With(ctx, zap.Stringer("remoteAddr", remoteAddr))

				err := s.DDOSProtector.Handle(ctx, conn)
				if err != nil {
					logger.FromCtx(ctx).Errorw("ddos protection", zap.Error(err))
					return
				}

				handler(ctx, conn)
			}()
		}
	}()

	return nil
}

// Shutdown shuts down the server by closing the listener and waiting for all
// the goroutines to finish executing.
func (s *Server) Shutdown() {
	if s.isClosed.CompareAndSwap(false, true) {
		s.l.Close()
	}
	s.wg.Wait()
}

// JoinErrors concatenates two error values into a single error.
func JoinErrors(err1, err2 error) error {
	switch {
	case err1 == nil:
		return err2
	case err2 == nil:
		return err1
	default:
		return errors.WithStack(stderrors.Join(err1, err2))
	}
}

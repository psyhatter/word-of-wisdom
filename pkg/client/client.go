package client

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"math"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/psyhatter/word-of-wisdom/pkg/models"
	"github.com/psyhatter/word-of-wisdom/pkg/puzzle"

	"github.com/pkg/errors"
)

// DDOSProtectedClient is a type that represents a client with distributed
// denial of service (DDoS) protection. It is used to establish a connection
// with a specified address and perform DDoS protection mechanisms.
type DDOSProtectedClient struct {
	// Address represents the address used to establish a connection.
	Address string

	// ConcurrencyLimit represents the maximum number of concurrent goroutines that
	// can be spawned to solve a puzzle.
	ConcurrencyLimit uint64

	wg  sync.WaitGroup
	cas atomic.Bool
}

// Connect establishes a connection with the specified address using the given
// context. The context can be used to cancel the net.DialContext operation
// only.
func (c *DDOSProtectedClient) Connect(ctx context.Context, handler func(context.Context, net.Conn) error) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", c.Address)
	if err != nil {
		return errors.Wrap(err, "dialing")
	}

	defer conn.Close()

	dec := json.NewDecoder(conn)

	var p models.Puzzle
	err = dec.Decode(&p)
	if err != nil {
		return errors.Wrap(err, "unmarshalling puzzle")
	}

	solution, err := c.solvePuzzle(p)
	if err != nil {
		return errors.Wrap(err, "solving puzzle")
	}

	err = json.NewEncoder(conn).Encode(solution)
	if err != nil {
		return errors.Wrap(err, "encoding solution")
	}

	var result models.VerificationResponse
	err = dec.Decode(&result)
	if err != nil {
		return errors.Wrap(err, "decoding response")
	}

	if !result.Success {
		return errors.New(result.ErrorMessage)
	}

	return errors.WithStack(handler(ctx, conn))
}

func (c *DDOSProtectedClient) solvePuzzle(p models.Puzzle) (solution models.Solution, err error) {
	if c.ConcurrencyLimit == 0 {
		c.ConcurrencyLimit = uint64(runtime.NumCPU())
	}

	c.wg.Add(int(c.ConcurrencyLimit))

	for i := uint64(0); i < c.ConcurrencyLimit; i++ {
		go func(i uint64) {
			defer c.wg.Done()

			buf := make([]byte, 0, unsafe.Sizeof(i))
			for ; i < math.MaxUint64-c.ConcurrencyLimit; i += c.ConcurrencyLimit {
				if !p.Deadline.IsZero() && time.Until(p.Deadline) <= 0 {
					if !c.cas.Load() && c.cas.CompareAndSwap(false, true) {
						err = errors.WithStack(ErrTimeout)
					}
					return
				}

				if c.cas.Load() {
					return
				}

				data := binary.BigEndian.AppendUint64(buf, i)
				if puzzle.CheckSolution(p, data) != nil {
					continue
				}

				if !c.cas.Load() && c.cas.CompareAndSwap(false, true) {
					solution = data
				}

				return
			}
		}(i)
	}

	c.wg.Wait()

	return
}

// ErrTimeout represents an error indicating a time-out occurred.
var ErrTimeout = errors.New("timeout")

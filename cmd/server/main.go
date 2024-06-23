package main

import (
	"context"
	"encoding/json"
	"flag"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/psyhatter/word-of-wisdom/internal/repo"
	"github.com/psyhatter/word-of-wisdom/internal/repo/inmemory"
	"github.com/psyhatter/word-of-wisdom/pkg/logger"
	"github.com/psyhatter/word-of-wisdom/pkg/server"

	"go.uber.org/zap"
)

var address = flag.String("address", ":8080", "server address for listening")
var timeout = flag.Duration("timeout", 10*time.Second, "time within which clients must send POW")
var nonceSize = flag.Uint64("nonceSize", 1024, "size of generated nonce")
var complexity = flag.Uint("complexity", 20, "complexity of tasks [0,255]")

func main() {
	flag.Parse()

	logger.SetLevel(zap.InfoLevel)

	ctx := context.Background()

	s := server.Server{
		Address: *address,
		DDOSProtector: server.DDOSProtector{
			Timeout:    *timeout,
			NonceSize:  int(*nonceSize),
			Complexity: uint8(*complexity),
		},
	}

	err := s.Listen(ctx, Service(inmemory.NewRepo()))
	if err != nil {
		logger.FromCtx(ctx).Fatalw("starting the server", zap.Error(err))
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	<-ctx.Done()
	s.Shutdown()
}

func Service(repo repo.Repo) func(ctx context.Context, conn net.Conn) {
	return func(ctx context.Context, conn net.Conn) {
		quote, err := repo.GetRandomQuote(ctx)
		if err != nil {
			logger.FromCtx(ctx).Errorw("getting random quote", zap.Error(err))
			return
		}

		err = json.NewEncoder(conn).Encode(quote)
		if err != nil {
			logger.FromCtx(ctx).Errorw("sending quote", zap.Error(err))
			return
		}
	}
}

package main

import (
	"context"
	"encoding/json"
	"flag"
	"net"
	"runtime"

	"github.com/pkg/errors"
	"github.com/psyhatter/word-of-wisdom/pkg/client"
	"github.com/psyhatter/word-of-wisdom/pkg/logger"

	"go.uber.org/zap"
)

var address = flag.String("address", "127.0.0.1:8080", "server address with ddos protection")
var workers = flag.Uint64("p", uint64(runtime.NumCPU()), "number of parallel workers to solve a POW task")

func main() {
	flag.Parse()

	logger.SetLevel(zap.InfoLevel)

	ctx := context.Background()

	c := client.DDOSProtectedClient{
		Address:          *address,
		ConcurrencyLimit: *workers,
	}

	err := c.Connect(ctx, Service)
	if err != nil {
		logger.FromCtx(ctx).Fatalw("connecting", zap.Error(err))
	}
}

func Service(ctx context.Context, conn net.Conn) error {
	var quote []byte
	err := json.NewDecoder(conn).Decode(&quote)
	if err != nil {
		return errors.Wrap(err, "decoding quote")
	}

	logger.FromCtx(ctx).Info(string(quote))

	return nil
}

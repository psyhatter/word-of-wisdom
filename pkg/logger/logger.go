package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetLevel sets the global logging level for the zap Logger. The provided level
// only enables logging levels that are greater or equal to it.
func SetLevel(lvl zapcore.LevelEnabler) {
	zap.ReplaceGlobals(
		zap.New(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.Lock(os.Stdout),
				lvl,
			),
		),
	)
}

type key struct{}

// ToCtx adds a SugaredLogger to the provided context. This logger can be
// retrieved using the FromCtx function.
func ToCtx(ctx context.Context, log *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, key{}, log)
}

// FromCtx retrieves the SugaredLogger from the provided context.
// If the SugaredLogger is not found in the context, it returns a default
// SugaredLogger.
func FromCtx(ctx context.Context) *zap.SugaredLogger {
	if s, _ := ctx.Value(key{}).(*zap.SugaredLogger); s != nil {
		return s
	}
	return zap.S()
}

// With adds the provided arguments to the given context and returns the modified context.
// Смотри (*zap.SugaredLogger) With
func With(ctx context.Context, args ...any) context.Context {
	return ToCtx(ctx, FromCtx(ctx).With(args...))
}

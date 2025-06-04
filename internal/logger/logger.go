package logger

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"strings"
)

type contextKey struct{}

// FromContextOrDiscard will return the slog.Logger from the context. If no slog.Logger
// is found, this will return a discard handler.
func FromContextOrDiscard(ctx context.Context) *slog.Logger {
	l, err := FromContext(ctx)
	if err != nil {
		l = DiscardingLogger()
	}
	return l
}

var ErrLoggerNotFoundInContext = errors.New("no pre-configured logger found")

func FromContext(ctx context.Context) (*slog.Logger, error) {
	l, ok := ctx.Value(contextKey{}).(*slog.Logger)
	if !ok {
		return nil, ErrLoggerNotFoundInContext
	}

	return l, nil
}

func NewContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

// New returns a structured logger given the provided inputs.
func New(level string, out io.Writer) (*slog.Logger, error) {
	levelUpper := strings.ToUpper(level)
	var loggerLevel slog.Level
	unrecognizedLevel := false
	switch levelUpper {
	case "DEBUG", "INFO", "WARN", "ERROR":
		err := loggerLevel.UnmarshalText([]byte(levelUpper))
		if err != nil {
			return nil, err
		}
	default: // The input log level was unrecognized
		unrecognizedLevel = true
		loggerLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{
		Level: loggerLevel,
	}))

	if unrecognizedLevel {
		logger.Warn("fallback log level was used because user-provided value was unrecognized", "provided", level)
	}

	return logger, nil
}

// MarshalJSON will produce a byte slice containing the result of a successful
// marshal of v, or the value "failedToMarshal" if the marshal attempt threw an
// error. This is designed to use with slog.
func MarshalJSON(v any) []byte {
	if b, err := json.Marshal(v); err == nil {
		return b
	}

	return []byte("failedToMarshal")
}

func DiscardingLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

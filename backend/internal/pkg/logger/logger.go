package logger

import (
	"context"
	"log/slog"
	"os"
)

// contextKey is unexported so nobody outside this package can clash with it.
type contextKey string

const (
	keyRequestID contextKey = "request_id"
	keyUserID    contextKey = "user_id"
)

// Setup initialises the global slog logger.
// Call this ONCE in main() before http.ListenAndServe.
func Setup(env string) {
	level := slog.LevelInfo
	if env == "dev" {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true, // adds "source":{"file":"...","line":N} to every record
	})

	slog.SetDefault(slog.New(handler))
}

// WithRequestID stores a request ID in the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, keyRequestID, requestID)
}

// WithUserID stores a user ID in the context.
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, keyUserID, userID)
}

// FromContext returns a *slog.Logger pre-enriched with
// request_id and user_id extracted from ctx.
// Always returns a usable logger even if context values are absent.
func FromContext(ctx context.Context) *slog.Logger {
	attrs := []any{}

	if rid, ok := ctx.Value(keyRequestID).(string); ok && rid != "" {
		attrs = append(attrs, slog.String("request_id", rid))
	}
	if uid, ok := ctx.Value(keyUserID).(int); ok && uid != 0 {
		attrs = append(attrs, slog.Int("user_id", uid))
	}

	return slog.Default().With(attrs...)
}

package logger

import (
	"context"
	"log/slog"
	"net/http"
)

type Logger = slog.Logger

type key string

const keyLogger key = "logger"

func NewDummy() *Logger {
	return slog.New(slog.NewJSONHandler(&nullWriter{}, nil))
}

type nullWriter struct {
}

func (nw *nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Inject
// create injects middleware
func Inject(l *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := l.With(
				slog.String("url", r.URL.String()),
				slog.String("method", r.Method),
				slog.Any("headers", r.Header),
			)
			r = r.WithContext(context.WithValue(r.Context(), keyLogger, l))
			next.ServeHTTP(w, r)
		})
	}
}

func FromContext(ctx context.Context) *Logger {
	return ctx.Value(keyLogger).(*Logger)
}

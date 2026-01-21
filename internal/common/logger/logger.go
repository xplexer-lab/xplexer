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
			l = l.WithGroup("request")
			r = r.WithContext(context.WithValue(r.Context(), keyLogger, l))
			next.ServeHTTP(w, r)
		})
	}
}

func Get(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(keyLogger).(*Logger); ok {
		return logger
	}
	return nil
}

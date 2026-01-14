package logger

import "log/slog"

func NewDummy() *slog.Logger {
	return slog.New(slog.NewJSONHandler(&nullWriter{}, nil))
}

type nullWriter struct {
}

func (nw *nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

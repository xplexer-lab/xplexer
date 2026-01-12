package api

import (
	"context"
	"log/slog"
)

type ContextType string

var (
	ContextKey = ContextType("api")
)

type Context interface {
	context.Context
	Logger() *slog.Logger
}

var _ Context = new(apiContext)

func GetCtx(ctx context.Context) Context {
	if cast, ok := ctx.Value(ContextKey).(*apiContext); ok {
		return cast
	}
	return nil
}

func WrapCtx(
	ctx context.Context,
	logger *slog.Logger,
) Context {
	if cast, ok := ctx.Value(ContextKey).(*apiContext); ok {
		return cast
	}

	newCtx := &apiContext{
		logger:  logger,
		Context: nil,
	}

	newCtx.Context = context.WithValue(ctx, ContextKey, newCtx)
	return newCtx
}

type apiContext struct {
	context.Context
	logger *slog.Logger
}

func (a apiContext) Logger() *slog.Logger {
	return a.logger
}

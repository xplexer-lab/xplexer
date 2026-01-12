package api

import "context"

type ContextType string

var (
	ContextKey = ContextType("api")
)

type Context interface {
	context.Context
	Logger() any // todo: select interface for logger
}

var _ Context = new(apiContext)

func WrapCtx(ctx context.Context) Context {
	if cast, ok := ctx.Value(ContextKey).(*apiContext); ok {
		return cast
	}

	newCtx := &apiContext{
		Context: nil,
	}

	newCtx.Context = context.WithValue(ctx, ContextKey, newCtx)
	return newCtx
}

type apiContext struct {
	context.Context
}

func (a apiContext) Logger() any {
	//TODO implement me
	panic("implement me")
}

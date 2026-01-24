package persistance

import (
	"context"
	"net/http"

	"github.com/xplexer-lab/xplexer/internal/common/errpack"
)

type txManagerKey[T any] struct{}

func Tx[Out, Repos any](
	ctx context.Context,
	handler func(context.Context, Repos) (Out, error),
	opts ...TxOpt,
) (Out, error) {
	var empty Out

	val := ctx.Value(txManagerKey[Repos]{})
	if val == nil {
		return empty, errpack.New(
			"TxManager missing in context: did you forget Inject middleware for this module?",
			errpack.Unreachable(),
		)
	}

	mgr, ok := val.(TxManager[Repos])
	if !ok {
		return empty, errpack.New(
			"TxManager type mismatch in context",
			errpack.Unreachable(),
		)
	}

	resAny, err := mgr.Do(
		ctx,
		func(txCtx context.Context, repos Repos) (any, error) {
			return handler(txCtx, repos)
		},
		opts...,
	)

	if err != nil {
		return empty, err
	}

	return resAny.(Out), nil
}

func Inject[Repos any](
	mgr TxManager[Repos],
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), txManagerKey[Repos]{}, mgr))
			next.ServeHTTP(w, r)
		})
	}
}

type TxCfg struct{}

type TxOpt func(*TxCfg)

type TxManager[Repos any] interface {
	Do(context.Context, func(context.Context, Repos) (any, error), ...TxOpt) (any, error)
}

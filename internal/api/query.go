package api

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/xplexer-lab/xplexer/internal/common/errpack"
)

var (
	_ Handler = new(queryHandler[any, any])
)

type (
	QueryOpt[In, Out any] func(*queryHandler[In, Out])

	queryHandler[In, Out any] struct {
		handlerCfg
		handle func(Context, In) (Out, error)
	}
)

func (qh *queryHandler[In, Out]) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// todo: inject logger into

	var in In
	// todo: bind In dto with request data
	var ctx = WrapCtx(r.Context())

	res, err := qh.handle(ctx, in)

	if err != nil {
		qh.handleError(rw, err)
		return
	}

	// todo: provide marshaller layer
	// marshaller can be dependent on Accept headers
	body, err := json.Marshal(res)
	if err != nil {
		// todo: log error
		qh.handleError(rw, errpack.Wrap(
			err,
			"failed to serialize json",
			errpack.WithDomain(),
		))
		return
	}

	if _, err = rw.Write(body); err != nil {
		// todo: log error
		return
	}

	// todo: parametrize ok status
	rw.WriteHeader(http.StatusOK)
}

func (qh *queryHandler[In, Out]) In() reflect.Type {
	return reflect.TypeFor[In]()
}

func (qh *queryHandler[In, Out]) Out() reflect.Type {
	return reflect.TypeFor[Out]()
}

func (qh *queryHandler[In, Out]) handleError(
	rw http.ResponseWriter,
	err error,
) {
	// todo: extract logger from context and log error
	rw.WriteHeader(http.StatusInternalServerError)
}

func Query[In, Out any](
	handle func(Context, In) (Out, error),
	opts ...QueryOpt[In, Out],
) Handler {
	res := &queryHandler[In, Out]{
		handle: handle,
	}

	for _, apply := range opts {
		apply(res)
	}

	return res
}

func WithQueryCommon[In, Out any](
	opts ...HandlerOpt,
) QueryOpt[In, Out] {
	return func(qh *queryHandler[In, Out]) {
		for _, apply := range opts {
			apply(&qh.handlerCfg)
		}
	}
}

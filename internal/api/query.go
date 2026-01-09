package api

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/xplexer-lab/xplexer/internal/common/errpack"
)

var (
	_ Handler[any] = new(queryHandler[any, any, any])
)

type (
	QueryOpt[In, Out, Ctx any] func(*queryHandler[In, Out, Ctx])

	queryHandler[In, Out, Ctx any] struct {
		handlerCfg
		handle func(*In, Ctx) (*Out, error)
	}
)

func (qh *queryHandler[In, Out, Ctx]) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	// todo: inject logger into

	var ctx Ctx
	// todo: extract ctx from context

	var in In
	// todo: bind In dto with request data

	res, err := qh.handle(&in, ctx)

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
			"failed to serizalize json",
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

func (qh *queryHandler[In, Out, Ctx]) In() reflect.Type {
	return reflect.TypeFor[In]()
}

func (qh *queryHandler[In, Out, Ctx]) Out() reflect.Type {
	return reflect.TypeFor[Out]()
}

func (qh *queryHandler[In, Out, Ctx]) handleError(
	rw http.ResponseWriter,
	err error,
) {
	// todo: extract logger from context and log error
	rw.WriteHeader(http.StatusInternalServerError)
}

func Query[In, Out, Ctx any](
	handle func(*In, Ctx) (*Out, error),
	opts ...QueryOpt[In, Out, Ctx],
) Handler[Ctx] {
	res := &queryHandler[In, Out, Ctx]{
		handle: handle,
	}

	for _, apply := range opts {
		apply(res)
	}

	return res
}

func WithQueryCommon[In, Out, Ctx any](
	opts ...HandlerOpt,
) QueryOpt[In, Out, Ctx] {
	return func(qh *queryHandler[In, Out, Ctx]) {
		for _, apply := range opts {
			apply(&qh.handlerCfg)
		}
	}
}

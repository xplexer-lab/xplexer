package restapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/creasty/defaults"
	validator "github.com/go-playground/validator/v10"

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
	var ctx = GetCtx(r.Context())
	ctx.Logger().Debug("debug starting query")

	in, err := qh.bind(r)
	if err != nil {
		qh.handleError(rw, err)
		return
	}

	res, err := qh.handle(ctx, in)

	if err != nil {
		ctx.Logger().Error(
			"request failed",
			slog.Any("err", err),
		)
		qh.handleError(rw, err)
		return
	}

	// Marshaller can be hidden behind abstract layer
	body, err := json.Marshal(res)
	if err != nil {
		ctx.Logger().ErrorContext(ctx, "failed to serialize json", slog.Any("err", err))
		qh.handleError(rw, errpack.Wrap(
			err,
			"failed to serialize json",
			errpack.WithDomain(),
		))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if _, err = rw.Write(body); err != nil {
		ctx.
			Logger().
			ErrorContext(
				ctx,
				"failed to write body",
				slog.Any("err", err),
			)
		return
	}
}

func (qh *queryHandler[In, Out]) In() reflect.Type {
	return reflect.TypeFor[In]()
}

func (qh *queryHandler[In, Out]) Out() reflect.Type {
	return reflect.TypeFor[Out]()
}

func (qh *queryHandler[In, Out]) bind(r *http.Request) (In, error) {
	var in In

	if err := defaults.Set(&in); err != nil {
		return in, err
	}

	if err := validator.New().Struct(in); err != nil {
		return in, err
	}

	return in, nil
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

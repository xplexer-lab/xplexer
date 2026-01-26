package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/creasty/defaults"
	validator "github.com/go-playground/validator/v10"

	"github.com/xplexer-lab/xplexer/internal/common/binder"
	"github.com/xplexer-lab/xplexer/internal/common/errpack"
)

var (
	_                Handler = new(operationHanlder[any, any])
	defaultBinder            = binder.Default()
	defaultValidator         = validator.New()
)

type (
	OperationOpt[In, Out any] func(*operationHanlder[In, Out])

	ErrorResponse struct {
		Error   string            `json:"error"`
		Details map[string]string `json:"details,omitempty"`
	}

	operationHanlder[In, Out any] struct {
		handlerCfg
		handle func(context.Context, In) (Out, error)
	}
)

func (qh *operationHanlder[In, Out]) Middlewares() Middlewares {
	return qh.middlewares
}

func (qh *operationHanlder[In, Out]) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var log = logger(ctx)
	log.Debug("debug starting query")

	in, err := qh.bind(r)
	if err != nil {
		qh.handleError(r, rw, err)

		return
	}

	res, err := qh.handle(ctx, in)

	if err != nil {
		qh.handleError(r, rw, err)
		return
	}

	body, err := json.Marshal(res)
	if err != nil {
		qh.handleError(r, rw, errpack.Wrap(
			err,
			"failed to serialize json",
			errpack.Domain(),
		))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(qh.getSuccessCode())
	if _, err = rw.Write(body); err != nil {
		log.
			ErrorContext(
				ctx,
				"failed to write body",
				slog.Any("err", err),
			)
		return
	}
}

func (qh *operationHanlder[In, Out]) In() reflect.Type {
	return reflect.TypeFor[In]()
}

func (qh *operationHanlder[In, Out]) Out() reflect.Type {
	return reflect.TypeFor[Out]()
}

func (qh *operationHanlder[In, Out]) bind(r *http.Request) (In, error) {
	var in In

	if err := defaults.Set(&in); err != nil {
		return in, err
	}

	if r.ContentLength > 0 && r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			return in, err
		}
	}

	if err := defaultBinder.Bind(r, &in); err != nil {
		return in, err
	}

	if casted, ok := any(&in).(Sanitizer); ok {
		casted.Sanitize()
	}

	if err := defaultValidator.Struct(in); err != nil {
		return in, err
	}

	return in, nil
}

func (qh *operationHanlder[In, Out]) handleError(
	r *http.Request,
	rw http.ResponseWriter,
	err error,
) {
	ctx := r.Context()
	status, shouldLogStack, response := qh.mapError(err)

	if shouldLogStack {
		slog.ErrorContext(ctx, "request failed with internal error", "err", err)
	} else {
		slog.InfoContext(ctx,
			"request failed with client error",
			slog.Any("err", err),
			slog.Int("status", status),
		)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)

	if encodeErr := json.NewEncoder(rw).Encode(response); encodeErr != nil {
		slog.ErrorContext(ctx,
			"failed to write error response",
			slog.Any("encodeErr", encodeErr),
		)
	}
}

func (qh *operationHanlder[In, Out]) mapError(err error) (int, bool, ErrorResponse) {
	// validation errors
	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		return http.StatusBadRequest, false, ErrorResponse{
			Error:   "Validation failed",
			Details: mapValidationErrors(valErrs),
		}
	}

	// json errors
	var jsonSyntaxErr *json.SyntaxError
	var jsonUnmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &jsonSyntaxErr) || errors.As(err, &jsonUnmarshalErr) {
		return http.StatusBadRequest, false, ErrorResponse{
			Error: "Invalid JSON format",
		}
	}

	// errpack errors
	var e *errpack.Error
	if errors.As(err, &e) {
		switch e.Type() {
		case errpack.TypeDomain:
			return http.StatusUnprocessableEntity, false, ErrorResponse{
				Error: e.Error(),
			}

		case errpack.TypeForbidden:
			return http.StatusForbidden, true, ErrorResponse{
				Error: "Access Denied",
			}

		case errpack.TypeInfra, errpack.TypeBootstrap:
			return http.StatusInternalServerError, true, ErrorResponse{
				Error: "Internal Server Error",
			}

		case errpack.TypeUnauthorized:
			return http.StatusUnauthorized, true, ErrorResponse{
				Error: "Unauthorized",
			}

		case errpack.TypeUnknown:
			fallthrough
		default:
			return qh.getErrorCodeFallback(), true, ErrorResponse{
				Error: "Internal Server Error",
			}
		}
	}

	return http.StatusInternalServerError, true, ErrorResponse{
		Error: "Internal Server Error",
	}
}

func Operation[In, Out any](
	handle func(context.Context, In) (Out, error),
	opts ...OperationOpt[In, Out],
) Handler {
	res := &operationHanlder[In, Out]{
		handle: handle,
	}

	for _, apply := range opts {
		apply(res)
	}

	return res
}

func WithOpCommon[In, Out any](
	opts ...HandlerOpt,
) OperationOpt[In, Out] {
	return func(qh *operationHanlder[In, Out]) {
		for _, apply := range opts {
			apply(&qh.handlerCfg)
		}
	}
}

func WithOpMiddlewares[In, Out any](
	middlewares ...Middleware,
) OperationOpt[In, Out] {
	return func(qh *operationHanlder[In, Out]) {
		qh.middlewares = append(qh.middlewares, middlewares...)
	}
}

func mapValidationErrors(verrs validator.ValidationErrors) map[string]string {
	result := make(map[string]string)
	for _, f := range verrs {
		result[f.Field()] = f.Tag()
	}
	return result
}

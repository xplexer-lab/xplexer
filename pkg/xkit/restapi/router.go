package restapi

import (
	"context"
	"log/slog"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/xplexer-lab/xplexer/pkg/xkit/errpack"
)

type (
	Router struct {
		routes      []route
		logger      *slog.Logger
		middlewares Middlewares
	}

	route struct {
		path    string
		method  string
		handler Handler
	}
)

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) SetLogger(logger *slog.Logger) *Router {
	r.logger = logger
	return r
}

func (r *Router) Get(path string, handler Handler) *Router {
	return r.Method(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler Handler) *Router {
	return r.Method(http.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler Handler) *Router {
	return r.Method(http.MethodPut, path, handler)
}

func (r *Router) Path(path string, handler Handler) *Router {
	return r.Method(http.MethodPatch, path, handler)
}

func (r *Router) Head(path string, handler Handler) *Router {
	return r.Method(http.MethodHead, path, handler)
}

func (r *Router) Options(path string, handler Handler) *Router {
	return r.Method(http.MethodOptions, path, handler)
}

// Adds global middlewares
func (r *Router) Use(middlewares ...Middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *Router) Method(
	method, path string,
	handler Handler,
) *Router {
	r.routes = append(r.routes, route{
		path:    path,
		method:  method,
		handler: handler,
	})
	return r
}

func (r *Router) BuildHandler() (http.Handler, error) {
	if r.logger == nil {
		return nil, errpack.New("logger is not provided", errpack.Bootstrap())
	}

	router := chi.NewRouter()
	router.Use(injectLogger(r.logger))
	router.Use(r.middlewares...)

	for _, ri := range r.routes {
		router.Method(ri.method, ri.path, ri.httpHandler())
	}

	return router, nil
}

// build handler and uses command scope middlewares
func (r route) httpHandler() http.Handler {
	var res http.Handler = r.handler

	mws := slices.Clone(r.handler.Middlewares())
	slices.Reverse(mws)

	for _, mw := range mws {
		res = mw(res)
	}

	return res
}

type loggerKeyType string

const loggerKey loggerKeyType = "logger"

func injectLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), loggerKey, logger)))
		})
	}
}

func logger(ctx context.Context) *slog.Logger {
	return ctx.Value(loggerKey).(*slog.Logger)
}

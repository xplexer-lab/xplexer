package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xplexer-lab/xplexer/internal/common/errpack"
)

type (
	Router struct {
		routes []route
		logger *slog.Logger
	}

	route struct {
		path    string
		method  string
		handler http.Handler
	}
)

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) SetLogger(logger *slog.Logger) *Router {
	r.logger = logger
	return r
}

func (r *Router) Get(path string, handler http.Handler) *Router {
	return r.Method(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler http.Handler) *Router {
	return r.Method(http.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler http.Handler) *Router {
	return r.Method(http.MethodPut, path, handler)
}

func (r *Router) Path(path string, handler http.Handler) *Router {
	return r.Method(http.MethodPatch, path, handler)
}

func (r *Router) Head(path string, handler http.Handler) *Router {
	return r.Method(http.MethodHead, path, handler)
}

func (r *Router) Options(path string, handler http.Handler) *Router {
	return r.Method(http.MethodOptions, path, handler)
}

func (r *Router) Method(
	method, path string,
	handler http.Handler,
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
		return nil, errpack.New("logger is not provided", errpack.WithBootstrap())
	}

	router := chi.NewRouter()
	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// todo: extract like a context injector
			handler.ServeHTTP(w, req.WithContext(
				WrapCtx(
					req.Context(),
					r.logger,
				),
			))
		})
	})

	for _, ri := range r.routes {
		router.Method(ri.method, ri.path, ri.handler)
	}

	return router, nil
}

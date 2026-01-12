package api

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
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

func (router *Router) SetLogger(logger *slog.Logger) *Router {
	router.logger = logger
	return router
}

func (router *Router) Get(path string, handler http.Handler) *Router {
	return router.method(http.MethodGet, path, handler)
}

func (router *Router) Post(path string, handler http.Handler) *Router {
	return router.method(http.MethodPost, path, handler)
}

func (router *Router) method(
	method, path string,
	handler http.Handler,
) *Router {
	router.routes = append(router.routes, route{
		path:    path,
		method:  method,
		handler: handler,
	})
	return router
}

func (router *Router) BuildHandler() (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// todo: extract like a context injector
			handler.ServeHTTP(w, r.WithContext(
				WrapCtx(
					r.Context(),
					router.logger,
				),
			))
		})
	})

	for _, rItem := range router.routes {
		r.Method(rItem.method, rItem.path, rItem.handler)
	}

	return r, nil
}

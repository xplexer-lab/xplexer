package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	Router struct {
		routes []route
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

func (b *Router) SetLogger() *Router {
	return b
}

func (b *Router) Get(path string, handler http.Handler) *Router {
	return b.method(http.MethodGet, path, handler)
}

func (b *Router) Post(path string, handler http.Handler) *Router {
	return b.method(http.MethodPost, path, handler)
}

func (b *Router) method(
	method, path string,
	handler http.Handler,
) *Router {
	b.routes = append(b.routes, route{
		path:    path,
		method:  method,
		handler: handler,
	})
	return b
}

func (b *Router) BuildHandler() (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// todo: extract like a context injector
			handler.ServeHTTP(w, r.WithContext(WrapCtx(r.Context())))
		})
	})

	for _, rItem := range b.routes {
		r.Method(rItem.method, rItem.path, rItem.handler)
	}

	return r, nil
}

package api

import (
	"errors"
	"net/http"
)

type (
	Builder[Ctx any] struct {
	}
)

func NewBuilder[Ctx any]() *Builder[Ctx] {
	return &Builder[Ctx]{}
}

func (b *Builder[Ctx]) WithErrorHandler() *Builder[Ctx] {
	return b
}

func (b *Builder[Ctx]) method(
	method, path string,
	handler http.Handler,
) *Builder[Ctx] {
	return b
}

func (b *Builder[Ctx]) Build() (http.Handler, error) {
	return nil, errors.New("not implemented")
}

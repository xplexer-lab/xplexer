package api

import "errors"

type (
	Api struct {
	}

	Builder struct {
	}
)

func (b *Builder) Build() (*Api, error) {
	return nil, errors.New("not implemented")
}

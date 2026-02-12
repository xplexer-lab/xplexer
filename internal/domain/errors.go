package domain

import "github.com/xplexer-lab/xplexer/pkg/xkit/errpack"

var (
	ErrEmptySlug      = errpack.New("empty slug values is not allowed", errpack.Domain())
	ErrNotImplemented = errpack.New("not implemented", errpack.Bootstrap())
)

package usecases

import (
	"context"

	"github.com/xplexer-lab/xplexer/internal/common/errpack"
	"github.com/xplexer-lab/xplexer/internal/common/restapi"
)

type ProjectCreateIn struct {
	Name string `json:"name"`
}

type ProjectCreateOut struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var ProjectCreate = restapi.Query(func(ctx context.Context, in ProjectCreateIn) (*ProjectCreateOut, error) {
	return nil, errpack.New("not implemented", errpack.Bootstrap())
})

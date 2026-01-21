package usecases

import (
	"context"
	"log/slog"

	"github.com/xplexer-lab/xplexer/internal/common/errpack"
	"github.com/xplexer-lab/xplexer/internal/common/logger"
	"github.com/xplexer-lab/xplexer/internal/common/restapi"
)

type ProjectCreateIn struct {
	Name         string `json:"name" validate:"required"`
	Token        string `header:"X-Token"`
	Content      string `header:"Content-type"`
	Page         int    `query:"page" default:"0"`
	ItemsPerPage int    `query:"items_per_age" default:"90"`
}

type ProjectCreateOut struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var projectCreate = restapi.Query(func(ctx context.Context, in ProjectCreateIn) (*ProjectCreateOut, error) {
	var log = logger.Get(ctx)

	log.Error("try execute query",
		slog.Any("in", in),
	)

	return nil, errpack.New("not implemented", errpack.Bootstrap())
})

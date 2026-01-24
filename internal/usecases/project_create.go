package usecases

import (
	"context"

	"github.com/xplexer-lab/xplexer/internal/common/errpack"
	"github.com/xplexer-lab/xplexer/internal/common/persistance"
	"github.com/xplexer-lab/xplexer/internal/common/restapi"
	"github.com/xplexer-lab/xplexer/internal/domain"
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

var projectCreate = restapi.Operation(func(ctx context.Context, in ProjectCreateIn) (*ProjectCreateOut, error) {
	return persistance.Tx(ctx, func(ctx context.Context, repos Repos) (*ProjectCreateOut, error) {
		project, err := domain.NewProject()

		if err != nil {
			return nil, err
		}

		if err := repos.Project.Insert(ctx, project); err != nil {
			return nil, err
		}

		return nil, errpack.New("not implemented", errpack.Bootstrap())
	})
})

package usecases

import (
	"context"
	"time"

	"github.com/xplexer-lab/xplexer/pkg/xkit/consistency"
	"github.com/xplexer-lab/xplexer/pkg/xkit/entity"
	"github.com/xplexer-lab/xplexer/pkg/xkit/restapi"
)

type ProjectGetIn struct {
	ID string `path:"id"`
}

type ProjectGetOut struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var projectGet = restapi.Operation(
	restapi.WithHandler(func(ctx context.Context, in ProjectGetIn) (*ProjectGetOut, error) {
		// todo: remove tx it is no needed here
		// todo: build common dtos
		return consistency.Tx(ctx, func(ctx context.Context, repos Repositories) (*ProjectGetOut, error) {
			// todo: move id parsing into the binder

			id, err := entity.ParseId(in.ID)

			if err != nil {
				return nil, err
			}

			if project, err := repos.Projects.FindOne(ctx, id); err != nil {
				return nil, err
			} else {
				return &ProjectGetOut{
					ID: project.Id().Hex(),
				}, nil
			}
		})
	}),
)

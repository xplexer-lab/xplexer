package usecases

import (
	"context"
	"time"

	"github.com/xplexer-lab/xplexer/internal/domain"
	"github.com/xplexer-lab/xplexer/pkg/xkit/consistency"
	"github.com/xplexer-lab/xplexer/pkg/xkit/restapi"
)

type ProjectCreateIn struct {
	Name string `json:"name" validate:"required"`
}

type ProjectCreateOut struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newProjectCreateOut(p *domain.Project) *ProjectCreateOut {
	return &ProjectCreateOut{
		ID:        p.Id().Hex(),
		Name:      p.Name(),
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
	}
}

var projectCreate = restapi.Operation(
	restapi.WithHandler(func(ctx context.Context, in ProjectCreateIn) (*ProjectCreateOut, error) {
		return consistency.Tx(ctx, func(ctx context.Context, repos Repositories) (*ProjectCreateOut, error) {
			project, err := domain.NewProject(in.Name)

			if err != nil {
				return nil, err
			}

			if err := repos.Projects.Insert(ctx, project); err != nil {
				return nil, err
			}

			return newProjectCreateOut(project), nil
		})
	}),
)

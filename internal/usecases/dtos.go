package usecases

import (
	"time"

	"github.com/xplexer-lab/xplexer/internal/domain"
)

type ProjectDto struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newProjectDto(p *domain.Project) *ProjectDto {
	return &ProjectDto{
		ID:        p.Id().Hex(),
		Name:      p.Name(),
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
	}
}

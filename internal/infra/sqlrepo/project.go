package sqlrepo

import (
	"github.com/xplexer-lab/xplexer/internal/domain"
	"github.com/xplexer-lab/xplexer/internal/usecases"
	"github.com/xplexer-lab/xplexer/pkg/xkit/entity"
	"github.com/xplexer-lab/xplexer/pkg/xkit/gormrepo"
	"gorm.io/gorm"
)

var (
	_ usecases.ProjectRepository = new(projectRepository)
)

type projectModel struct {
	gormrepo.Model
	Name string
}

type projectRepository struct {
	*gormrepo.Repository[*domain.Project, domain.ProjectState, projectModel]
}

func newProjectRepository(tx *gorm.DB) *projectRepository {
	base := gormrepo.New[*domain.Project, domain.ProjectState, projectModel](
		tx,
		func() *domain.Project {
			return &domain.Project{
				Entity: &entity.Entity{},
			}
		},
		func(ps domain.ProjectState) (*projectModel, error) {
			pm := &projectModel{}

			if err := pm.LoadState(ps.State); err != nil {
				return nil, err
			}

			pm.Name = ps.Name

			return pm, nil
		},
		func(pm *projectModel) (domain.ProjectState, error) {
			var state domain.ProjectState
			var err error

			if state.State, err = pm.ToState(); err != nil {
				return state, nil
			}

			state.Name = pm.Name

			return state, nil
		},
	)

	return &projectRepository{Repository: base}
}

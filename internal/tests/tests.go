package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xplexer-lab/xplexer/internal/domain"
	"github.com/xplexer-lab/xplexer/internal/usecases"
	"github.com/xplexer-lab/xplexer/pkg/xkit/consistency"
	"github.com/xplexer-lab/xplexer/pkg/xkit/entity"
)

type (
	Tests struct {
		txm     TxManager
		handler http.Handler
	}

	// aliases
	TxManager = consistency.TxManager[usecases.Repositories]
)

// Acceptance tests
func New(
	txm TxManager,
	handler http.Handler,
) *Tests {
	return &Tests{
		txm:     txm,
		handler: handler,
	}
}

func (tt *Tests) Run(t *testing.T) {
	t.Run("Project", tt.testProjects)
}

func (tt *Tests) test(t testing.TB) *httpexpect.Expect {
	server := httptest.NewServer(tt.handler)
	return httpexpect.Default(t, server.URL)
}

func (tt *Tests) tx(t *testing.T, cb func(context.Context, usecases.Repositories) error) error {
	_, err := tt.txm.Do(t.Context(), func(ctx context.Context, r usecases.Repositories) (any, error) {
		return nil, cb(ctx, r)
	})

	return err
}

func (tt *Tests) testProjects(t *testing.T) {
	t.Helper()

	t.Run("GET /core/projects/{id}", func(t *testing.T) {
		t.Run("project is avaliable by id", func(t *testing.T) {
			project, err := domain.NewProject("my-project")
			require.NoError(t, err)

			_, err = tt.txm.Do(t.Context(), func(ctx context.Context, r usecases.Repositories) (any, error) {
				return nil, r.Projects.Insert(t.Context(), project)
			})

			require.NoError(t, err)

			t.Cleanup(func() {
				err := tt.txm.Repos().Projects.Delete(context.TODO(), project.Id())
				assert.NoError(t, err)
			})

			var res usecases.ProjectDto

			tt.test(t).
				GET("/core/projects/{id}", project.Id().Hex()).
				Expect().
				Status(http.StatusOK).
				JSON().
				Decode(&res)

			assert.Equal(t, "my-project", res.Name)
		})

		t.Run("returns 422 if project does not exist", func(t *testing.T) {
			tt.test(t).
				GET("/core/projects/{id}", entity.NewId().Hex()).
				Expect().
				Status(http.StatusUnprocessableEntity).
				JSON()
		})
	})

	t.Run("GET /core/projects", func(t *testing.T) {
		t.Run("returns list of existend projects", func(t *testing.T) {
			p1 := domain.MustNewProject("p1")
			p2 := domain.MustNewProject("p2")
			p3 := domain.MustNewProject("p3")

			tt.tx(t, func(ctx context.Context, r usecases.Repositories) error {
				return errors.Join(
					r.Projects.Insert(ctx, p1),
					r.Projects.Insert(ctx, p2),
					r.Projects.Insert(ctx, p3),
				)
			})

			t.Cleanup(func() {
				ctx := context.TODO()

				err := errors.Join(
					tt.txm.Repos().Projects.Delete(ctx, p1.Id()),
					tt.txm.Repos().Projects.Delete(ctx, p2.Id()),
					tt.txm.Repos().Projects.Delete(ctx, p3.Id()),
				)

				assert.NoError(t, err)
			})

			var dtos []usecases.ProjectDto

			tt.test(t).
				GET("/core/projects").
				Expect().
				Status(http.StatusOK).
				JSON().
				Decode(&dtos)

			assert.Equal(t, 3, len(dtos))
			assert.ElementsMatch(
				t,
				[]string{p1.Name(), p2.Name(), p3.Name()},
				lo.Map(dtos, func(dto usecases.ProjectDto, _ int) string {
					return dto.Name
				}),
			)
		})

		t.Run("returns empty list of projects", func(t *testing.T) {
			var dtos []usecases.ProjectDto

			tt.test(t).
				GET("/core/projects").
				Expect().
				Status(http.StatusOK).
				JSON().Decode(&dtos)

			assert.Len(t, dtos, 0)
		})
	})
}

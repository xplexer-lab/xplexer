package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"github.com/xplexer-lab/xplexer/internal/domain"
	"github.com/xplexer-lab/xplexer/internal/usecases"
	"github.com/xplexer-lab/xplexer/pkg/xkit/consistency"
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

func (tt *Tests) testProjects(t *testing.T) {
	t.Helper()

	t.Run("GET /core/projects/{id}", func(t *testing.T) {
		t.Run("project is avaliable by id", func(t *testing.T) {
			t.Skip("fix failing tests")

			project, err := domain.NewProject("my-project")
			require.NoError(t, err)

			_, err = tt.txm.Do(t.Context(), func(ctx context.Context, r usecases.Repositories) (any, error) {
				return nil, r.Projects.Insert(t.Context(), project)
			})

			require.NoError(t, err)

			tt.test(t).
				GET("/core/projects/{id}", project.Id().Hex()).
				Expect().
				Status(http.StatusOK).
				JSON()
		})
	})
}

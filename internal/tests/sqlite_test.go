package tests_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xplexer-lab/xplexer/internal/infra/sqlrepo"
	"github.com/xplexer-lab/xplexer/internal/tests"
	"github.com/xplexer-lab/xplexer/internal/usecases"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSqlite(t *testing.T) {
	const db = "xplexer-test.db"
	_ = os.Remove(db)

	sqlite.Open(db)

	repo, err := sqlrepo.NewTxManager(sqlite.Open(db), &gorm.Config{})
	require.NoError(t, err)

	hander, err := usecases.BuildRouter(repo).BuildHandler()
	tests.New(repo, hander).Run(t)
}

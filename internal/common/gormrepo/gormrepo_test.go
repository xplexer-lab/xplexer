package gormrepo_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/xplexer-lab/xplexer/internal/common/entitytest"
	"github.com/xplexer-lab/xplexer/internal/common/gormrepo"
	"github.com/xplexer-lab/xplexer/internal/common/testutils"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type model struct {
	gormrepo.Model
	Name string
}

func buildGormRepo(db *gorm.DB) entitytest.Repository {
	return gormrepo.New(
		db,
		entitytest.EmptyDummy,
		func(ds entitytest.State) (*model, error) {
			var model model
			model.Name = ds.Name
			return &model, model.LoadState(ds.State)
		},
		func(dm *model) (entitytest.State, error) {
			s, err := dm.ToState()

			return entitytest.State{
				State: s,
				Name:  dm.Name,
			}, err
		},
	)
}

func TestGormRepoPostgress(t *testing.T) {
	testutils.SkipIntegral(t)

	ctx := t.Context()

	dbName := "xplexer"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(t.Context(),
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
		testcontainers.WithExposedPorts(),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	host, err := postgresContainer.Host(ctx)
	port, err := postgresContainer.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		host,
		port.Port(),
		dbUser,
		dbPassword,
		dbName,
	)

	db, err := gorm.Open(gormpg.Open(dsn))
	require.NoError(t, err)

	err = db.AutoMigrate(&model{})
	require.NoError(t, err)

	entitytest.New(buildGormRepo(db)).Run(t)
}

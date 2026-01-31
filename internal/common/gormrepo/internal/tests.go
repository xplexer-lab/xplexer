package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xplexer-lab/xplexer/internal/common/entity"
	"github.com/xplexer-lab/xplexer/internal/common/errpack"
	"github.com/xplexer-lab/xplexer/internal/common/gormrepo"
	"gorm.io/gorm"
)

type Tests struct {
}

var (
	_ entity.Aggregate[dummyState] = new(dummy)
)

type dummy struct {
	*entity.Entity
	Name string
}

func newDummy() *dummy {
	return &dummy{
		Name:   "",
		Entity: entity.New(),
	}
}

func (du *dummy) ToState() dummyState {
	return dummyState{
		State: du.Entity.ToState(),
		Name:  du.Name,
	}
}

func (du *dummy) Load(s dummyState) error {
	du.Name = s.Name

	return du.Entity.Load(s.State)
}

type dummyState struct {
	entity.State
	Name string
}

type dummyModel struct {
	gormrepo.Model
	Name string
}

type dummyDepository = *gormrepo.Repository[*dummy, dummyState, dummyModel]

func newRepo(db *gorm.DB) dummyDepository {
	return gormrepo.New[*dummy, dummyState, dummyModel](
		db,
		newDummy,
		func(ds dummyState) (*dummyModel, error) {
			return nil, errpack.New("to state is not implemented")
		},
		func(dm *dummyModel) (dummyState, error) {
			return dummyState{}, errpack.New("to mode is not implemented")
		},
	)
}

func NewTests() *Tests {
	return &Tests{}
}

func (tt *Tests) Run(t *testing.T, db *gorm.DB) {
	t.Helper()

	err := db.AutoMigrate(&dummyModel{})
	require.NoError(t, err, "failed to migrate")

	repo := newRepo(db)
	require.NotNil(t, repo, "repository is created")

	t.Run("FindOne", func(t *testing.T) {

		t.Run("returns not found error", func(t *testing.T) {
			entry, err := repo.FindOne(t.Context(), entity.NewId())
			t.Logf("error found: %#v", err)
			assert.Nil(t, entry, "requires entity to return empty ")
			assert.ErrorIs(t, err, entity.ErrNotFound)
		})
	})
}

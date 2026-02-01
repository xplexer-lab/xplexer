// entitytest
// acceptance tests for abstract repository
package entitytest

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xplexer-lab/xplexer/internal/common/entity"
)

type Test struct {
	repo Repository
}

var (
	_ entity.Aggregate[State] = new(Entity)
)

type Entity struct {
	*entity.Entity
	Name     string
	validate func() error
}

func (d *Entity) Validate() error {
	if d.validate == nil {
		return nil
	}
	return d.validate()
}

func NewDummy(name string) *Entity {
	return &Entity{
		Entity:   entity.New(),
		Name:     name,
		validate: skipValidation,
	}
}

func EmptyDummy() *Entity {
	return &Entity{
		Entity:   &entity.Entity{},
		validate: skipValidation,
	}
}

func skipValidation() error {
	return nil
}

func (du *Entity) ToState() State {
	return State{
		State: du.Entity.ToState(),
		Name:  du.Name,
	}
}

func (du *Entity) Load(s State) error {
	du.Name = s.Name

	return du.Entity.Load(s.State)
}

type State struct {
	entity.State
	Name string
}

// Repository represents abstract repository
type Repository interface {
	entity.OneFinder[*Entity, State]
	entity.Updater[*Entity, State]
	entity.Saver[*Entity, State]
	entity.Inserter[*Entity, State]
	entity.Deleter[*Entity, State]
}

func New(
	repo Repository,
) *Test {
	return &Test{
		repo: repo,
	}
}

func (tt *Test) Run(t *testing.T) {
	t.Helper()
	repo := tt.repo

	require.NotNil(t, repo, "expected repository to be defined")

	t.Run("FindOne", func(t *testing.T) {
		t.Run("returns not found error", func(t *testing.T) {
			entry, err := repo.FindOne(t.Context(), entity.NewId())
			assert.Nil(t, entry, "requires entity to return empty ")
			assert.ErrorIs(t, err, entity.ErrNotFound)
		})

		t.Run("returns previously inserted item", func(t *testing.T) {
			entry := NewDummy("returns previously inserted")
			err := repo.Insert(t.Context(), entry)
			require.NoError(t, err)

			assert.True(t, entry.Name == "returns previously inserted")
			restored, err := repo.FindOne(t.Context(), entry.Id())

			assert.NoError(t, err)
			assert.NotNil(t, restored)

			assert.Equal(t, entry.Name, restored.Name)
			assert.Equal(t, entry.Id(), restored.Id())
			assert.Equal(t, entry.CreatedAt(), restored.CreatedAt())
			assert.Equal(t, entry.UpdatedAt(), restored.UpdatedAt())
			assert.Equal(t, entry.Version(), restored.Version())
		})
	})

	t.Run("Insert", func(t *testing.T) {
		t.Run("inserts into database", func(t *testing.T) {
			entry := NewDummy("hello")
			err := repo.Insert(t.Context(), entry)
			require.NoError(t, err)

			assert.True(t, entry.Name == "hello")
			restored, err := repo.FindOne(t.Context(), entry.Id())

			assert.NoError(t, err)
			assert.NotNil(t, restored)
			assert.Equal(t, entry.Name, restored.Name)
			assert.Equal(t, entry.ToState(), restored.ToState())
		})

		t.Run("invokes validate method if one exists", func(t *testing.T) {
			var expectedErr = errors.New("artifitial validate error")
			entry := NewDummy("hello")
			entry.validate = func() error {
				return expectedErr
			}
			err := repo.Insert(t.Context(), entry)
			require.ErrorIs(t, err, expectedErr)
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("returns not found error", func(t *testing.T) {
			err := repo.Update(t.Context(), entity.NewId(), func(_ *Entity) error {
				return nil
			})

			assert.ErrorIs(t, err, entity.ErrNotFound)
			assert.True(t, entity.ErrNotFound.EqHead(err))
		})

		t.Run("updates entry", func(t *testing.T) {
			entry := NewDummy("start")
			assert.True(t, entry.Version() == 0, "has default version")
			err := repo.Insert(t.Context(), entry)
			assert.True(t, entry.Version() == 1, "bumps version on insert")
			require.NoError(t, err)

			err = repo.Update(t.Context(), entry.Id(), func(e *Entity) error {
				e.Name = "updated"
				return nil
			})
			require.NoError(t, err)
			updated, err := repo.FindOne(t.Context(), entry.Id())
			require.NoError(t, err)

			assert.Equal(t, "updated", updated.Name)
			assert.Equal(t, 2, updated.Version())
		})

	})

	t.Run("Save", func(t *testing.T) {
		t.Run("bumps version each time", func(t *testing.T) {
			entry := NewDummy("0")
			assert.True(t, entry.Version() == 0, "has default version")

			err := repo.Insert(t.Context(), entry)
			assert.NoError(t, err, "saves")

			for range 3 {
				prevVersion := entry.Version()
				err := repo.Save(t.Context(), entry)
				assert.NoError(t, err)
				assert.Equal(t, prevVersion+1, entry.Version(), "bumpts version each time")
			}

		})
	})
}

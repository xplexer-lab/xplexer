package entity_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xplexer-lab/xplexer/internal/common/entity"
)

type dummy struct {
	*entity.Entity
	title string
	desc  string
}

func newDummy() *dummy {
	return &dummy{
		Entity: entity.New(),
		title:  "hello",
		desc:   "desc",
	}
}

func TestEntity(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		t.Run("delete marks entity as deleted", func(t *testing.T) {
			d := newDummy()
			assert.False(t, d.Deleted())
			d.Delete()
			assert.True(t, d.Deleted())
			deletedAt := d.DeletedAt()
			assert.NotNil(t, deletedAt)
			d.Delete()
			assert.Equal(t, deletedAt, d.DeletedAt())
		})
	})
}

func TestId(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		id := entity.NewId()
		str := strings.ReplaceAll(id.Hex(), "-", "")
		assert.True(t, len(str) == 32, "16 bytes id is created")
		var zeroId entity.Id
		assert.True(t, zeroId != id, "creates non zero struct")
	})
}

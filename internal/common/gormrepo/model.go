package gormrepo

import (
	"time"

	"github.com/xplexer-lab/xplexer/internal/common/entity"
	"gorm.io/gorm"
)

var (
	_ entity.Versioner = new(Model)
)

type Model struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	V         int            `gorm:"column:__v;default:1"`
}

func (mo *Model) Version() int {
	return mo.V
}

func (mo *Model) SetVersion(v int) {
	mo.V = v
}

func (m *Model) ToState() (entity.State, error) {
	id, err := entity.ParseId(m.ID)

	if err != nil {
		return entity.State{}, err
	}

	return entity.State{
		Id:        id,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeleteAt:  &m.DeletedAt.Time,
		Version:   m.V,
	}, nil
}

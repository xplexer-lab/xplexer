package gormrepo

import (
	"context"

	"github.com/xplexer-lab/xplexer/internal/common/entity"
	"gorm.io/gorm"
)

var (
	_ entity.Inserter[entity.Aggregate[any], any]  = new(Repository[entity.Aggregate[any], any, any])
	_ entity.OneFinder[entity.Aggregate[any], any] = new(Repository[entity.Aggregate[any], any, any])
	// todo: implement other parts of repository
)

// Creates empty entity
type Factory[E any] func() E

// State <-> Model mapping and vice versa
type ToModelFn[S any, M any] func(S) *M
type ToStateFn[S any, M any] func(*M) S

type Repository[E entity.Aggregate[S], S any, M any] struct {
	db      *gorm.DB
	factory Factory[E]
	toModel ToModelFn[S, M]
	toState ToStateFn[S, M]
}

func New[E entity.Aggregate[S], S any, M any](
	db *gorm.DB,
	factory Factory[E],
	toModel ToModelFn[S, M],
	toState ToStateFn[S, M],
) *Repository[E, S, M] {
	return &Repository[E, S, M]{
		db:      db,
		factory: factory,
		toModel: toModel,
		toState: toState,
	}
}

func (r *Repository[E, S, M]) Insert(ctx context.Context, ent E) error {
	state := ent.ToState()
	model := r.toModel(state)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *Repository[E, S, M]) FindOne(ctx context.Context, id entity.Id) (E, error) {
	var model M
	var empty E

	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return empty, err
	}

	state := r.toState(&model)

	ent := r.factory()
	if err := ent.Load(state); err != nil {
		return empty, err
	}

	return ent, nil
}

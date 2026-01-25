package gormrepo

import (
	"context"

	"github.com/xplexer-lab/xplexer/internal/common/entity"
	"gorm.io/gorm"
)

var (
	_ entity.Inserter[entity.Aggregate[any], any]  = new(Repository[entity.Aggregate[any], any, Model])
	_ entity.OneFinder[entity.Aggregate[any], any] = new(Repository[entity.Aggregate[any], any, Model])
	_ entity.Updater[entity.Aggregate[any], any]   = new(Repository[entity.Aggregate[any], any, Model])
	_ entity.Deleter[entity.Aggregate[any], any]   = new(Repository[entity.Aggregate[any], any, Model])
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

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// backward propagation
	return ent.Load(r.toState(model))
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

func (r *Repository[E, S, M]) Update(ctx context.Context, id entity.Id, update func(E) error) error {
	ent, err := r.FindOne(ctx, id)

	if err != nil {
		return err
	}

	if err := update(ent); err != nil {
		return err
	}

	version := ent.Version()
	ent.SetVersion(version + 1)
	newModel := r.toModel(ent.ToState())
	ent.SetVersion(version)

	result := r.db.WithContext(ctx).
		Model(newModel).
		Omit("id", "created_at").
		Where("id = ? AND __v = ?", id.Hex(), version).
		Updates(newModel)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return entity.ErrOptimisticLock
	}

	return nil
}

func (r *Repository[E, S, M]) Delete(ctx context.Context, id entity.Id) error {
	var model M
	return r.db.WithContext(ctx).
		Where("id = ?", id.Hex()).
		Delete(&model).Error
}


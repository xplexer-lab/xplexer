package entity

import (
	"context"
	"time"

	"github.com/xplexer-lab/xplexer/pkg/xkit/bus"
	"github.com/xplexer-lab/xplexer/pkg/xkit/errpack"
	"github.com/xplexer-lab/xplexer/pkg/xkit/query"
)

var (
	ErrNotFound       = errpack.New("entity not found", errpack.Domain())
	ErrOptimisticLock = errpack.New("err optimistic lock", errpack.Infra())
)

type Dumper[State any] interface {
	ToState() State
}

type Loader[State any] interface {
	Load(State) error
}

type EventProvider interface {
	PopEvents() []bus.AnyEnvelope
}

type Versioner interface {
	Version() int
	SetVersion(v int)
}

type Aggregate[State any] interface {
	Dumper[State]
	Loader[State]
	EventProvider
	Versioner
	Entitier
}

type Entitier interface {
	Id() Id
	CreatedAt() time.Time
	UpdatedAt() time.Time
	SetUpdatedAt(time.Time)
}

type OneFinder[T Aggregate[S], S any] interface {
	FindOne(ctx context.Context, id Id) (T, error)
}

type Finder[T Aggregate[S], S any] interface {
	Find(ctx context.Context, cond query.Condition) ([]T, error)
}

type Updater[T Aggregate[S], S any] interface {
	Update(ctx context.Context, id Id, update func(T) error) error
}

type Saver[T Aggregate[S], S any] interface {
	Save(ctx context.Context, entity T) error
}

type Inserter[T Aggregate[S], S any] interface {
	Insert(ctx context.Context, entity T) error
}

type Deleter[T Aggregate[S], S any] interface {
	Delete(ctx context.Context, id Id) error
}

type Repository[T Aggregate[S], S any] interface {
	OneFinder[T, S]
	Finder[T, S]
	Updater[T, S]
	Saver[T, S]
	Inserter[T, S]
	Deleter[T, S]
}

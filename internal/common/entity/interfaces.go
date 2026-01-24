package entity

import (
	"context"

	"github.com/xplexer-lab/xplexer/internal/common/errpack"
	"github.com/xplexer-lab/xplexer/internal/common/query"
	"google.golang.org/protobuf/proto"
)

var (
	ErrNotFound = errpack.New("entity not found", errpack.Domain())
)

type Dumper[State any] interface {
	ToState() State
}

type Loader[State any] interface {
	Load(State) error
}

type EventProvider interface {
	PopEvents() []proto.Message
}

type Aggregate[State any] interface {
	Dumper[State]
	Loader[State]
	EventProvider
}

type OneFinder[T Aggregate[S], S any] interface {
	FindOne(ctx context.Context, id Id) (T, error)
}

type Finder[T Aggregate[S], S any] interface {
	Find(ctx context.Context, cond query.Condition) ([]T, error)
}

type Updater[T Aggregate[S], S any] interface {
	Update(ctx context.Context, id Id, update func(*T) error) error
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
	Inserter[T, S]
	Deleter[T, S]
}

package bus

import (
	"context"
	"io"
)

type Handler[E AnyEnvelope] interface {
	Handle(msg E) error
}

type HandleFn[E AnyEnvelope] func(msg E) error

func (fn HandleFn[E]) Handle(msg E) error {
	return fn(msg)
}

type AnyHandler Handler[AnyEnvelope]

type Cancel func()

type Publisher interface {
	Publish(context.Context, AnyEnvelope) error
}

type Subscriber interface {
	Subscribe(context.Context, AnyEnvelope, AnyHandler) Cancel
}

type Bus interface {
	Publisher
	Subscriber
	io.Closer
}

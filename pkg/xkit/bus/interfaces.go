package bus

import (
	"context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"io"
)

type Handler[P proto.Message] interface {
	Handle(ctx context.Context, msg Envelope[P]) error
	Name() string
}

type AnyHandler interface {
	Handler[proto.Message]
	getProtoMessage() protoreflect.Message
}

func NewAnyHandlerWithHandler[P proto.Message](handler Handler[P]) AnyHandler {
	return NewAnyHandler(handler.Name(), handler.Handle)
}

func NewAnyHandler[P proto.Message](name string, handleFn func(context.Context, Envelope[P]) error) AnyHandler {
	var payload P
	var _ = payload

	return &anyHandlerFn{
		name:         name,
		protoMessage: payload.ProtoReflect(),
		fn: func(ctx context.Context, msg Envelope[proto.Message]) error {
			// todo: cast generic message onto typed
			return handleFn(ctx, NewEnvelope(payload))
		},
	}
}

var _ AnyHandler = new(anyHandlerFn)

type anyHandlerFn struct {
	name         string
	protoMessage protoreflect.Message
	fn           func(ctx context.Context, msg Envelope[proto.Message]) error
}

func (a *anyHandlerFn) Handle(ctx context.Context, msg Envelope[proto.Message]) error {
	return a.fn(ctx, msg)
}

func (a *anyHandlerFn) Name() string {
	return a.name
}

func (a *anyHandlerFn) getProtoMessage() protoreflect.Message {
	return a.protoMessage
}

type Cancel func()

type Publisher interface {
	Publish(context.Context, AnyEnvelope) error
}

type Subscriber interface {
	Subscribe(context.Context, AnyHandler) (Cancel, error)
}

type Bus interface {
	Publisher
	Subscriber
	io.Closer
}

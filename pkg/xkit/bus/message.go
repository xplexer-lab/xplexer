package bus

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

var (
	_ AnyEnvelope = new(Envelope[proto.Message])
)

func NewEnvelope[P proto.Message](payload P, opts ...MessageOpt) *Envelope[P] {
	msg := Envelope[P]{
		payload: payload,
		metadata: Metadata{
			ID:        uuid.New().String(),
			CreatedAt: time.Now(),
		},
	}

	for _, opt := range opts {
		opt(&msg.metadata)
	}

	return &msg
}

type AnyEnvelope interface {
	ID() string
	Metadata() Metadata
	ProtoPayload() proto.Message
}

type Envelope[P proto.Message] struct {
	payload  P
	metadata Metadata
}

type Metadata struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	Attempt   int        `json:"attempt"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (e *Envelope[P]) ID() string {
	return e.metadata.ID
}

func (e *Envelope[P]) Metadata() Metadata {
	return e.metadata
}

func (e *Envelope[P]) ProtoPayload() proto.Message {
	return e.payload
}

func (e *Envelope[P]) Payload() P {
	return e.payload
}

func (e *Envelope[P]) IncrAttempts() {
	e.metadata.Attempt++
}

type MessageOpt func(*Metadata)

func WithMetadata(metadata Metadata) MessageOpt {
	return func(m *Metadata) {
		*m = metadata
	}
}

func WithExpiration(expiresAt time.Time) MessageOpt {
	return func(m *Metadata) {
		m.ExpiresAt = &expiresAt
	}
}

func WithCreatedAt(createdAt time.Time) MessageOpt {
	return func(m *Metadata) {
		m.CreatedAt = createdAt
	}
}

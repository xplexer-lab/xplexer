package bus

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

var (
	_ AnyEnvelope = new(Envelope[proto.Message])
)

func NewEnvelope[P proto.Message](payload P, opts ...EnvelopeOpt) *Envelope[P] {
	msg := Envelope[P]{
		payload: payload,
		metadata: Metadata{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
		},
	}

	for _, opt := range opts {
		opt(&msg.metadata)
	}

	return &msg
}

type AnyEnvelope interface {
	ID() uuid.UUID
	Metadata() Metadata
	ProtoPayload() proto.Message
}

type Envelope[P proto.Message] struct {
	payload  P
	metadata Metadata
}

type Metadata struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	Attempt   int        `json:"attempt"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (e *Envelope[P]) ID() uuid.UUID {
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

type EnvelopeOpt func(*Metadata)

func WithMetadata(metadata Metadata) EnvelopeOpt {
	return func(m *Metadata) {
		*m = metadata
	}
}

func WithExpiration(expiresAt time.Time) EnvelopeOpt {
	return func(m *Metadata) {
		m.ExpiresAt = &expiresAt
	}
}

func WithCreatedAt(createdAt time.Time) EnvelopeOpt {
	return func(m *Metadata) {
		m.CreatedAt = createdAt
	}
}

func WithId(id uuid.UUID) EnvelopeOpt {
	return func(m *Metadata) {
		m.ID = id
	}
}

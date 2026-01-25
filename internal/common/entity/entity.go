package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/xplexer-lab/xplexer/internal/common/bus"
	"github.com/xplexer-lab/xplexer/internal/common/errpack"
)

var (
	// Base entity has to implement Aggregate
	_ Aggregate[State] = new(Entity)
	_ Versioner        = new(Entity)
	_ EventProvider    = new(Entity)
	_ Loader[State]    = new(Entity)
	_ Dumper[State]    = new(Entity)
)

type Id uuid.UUID

func NewId() Id {
	return Id(uuid.New())
}

func (id Id) Hex() string {
	return uuid.UUID(id).String()
}

type Entity struct {
	id        Id
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
	version   int

	events []bus.AnyEnvelope
}

func (e *Entity) Version() int {
	return e.version
}

func (e *Entity) SetVersion(v int) {
	e.version = v
}

type State struct {
	Id        Id
	CreatedAt time.Time
	UpdatedAt time.Time
	DeleteAt  *time.Time
	Version   int
}

type Opt func(*Entity)

func New(opts ...Opt) *Entity {
	return &Entity{
		id:        NewId(),
		createdAt: time.Now(),
		updatedAt: time.Now(),
		deletedAt: nil,
		version:   0,
	}
}

func (e *Entity) Id() Id {
	return e.id
}

func (e *Entity) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *Entity) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Entity) DeletedAt() *time.Time {
	return e.deletedAt
}

func (e *Entity) Deleted() bool {
	return e.deletedAt != nil
}

func (e *Entity) DeleteAt(ts time.Time) {
	if e.deletedAt != nil {
		return
	}

	e.deletedAt = &ts
}

func (e *Entity) Delete() {
	e.DeleteAt(time.Now())
}

func (e *Entity) ToState() State {
	return State{
		Id:        e.id,
		CreatedAt: e.createdAt,
		UpdatedAt: e.updatedAt,
		DeleteAt:  e.deletedAt,
		Version:   e.version,
	}
}

func (e *Entity) RecordEvent(evt bus.AnyEnvelope) {
	e.events = append(e.events, evt)
}

func (e *Entity) PopEvents() []bus.AnyEnvelope {
	ret := e.events
	e.events = nil
	return ret
}

func (e *Entity) Load(s State) error {
	if len(e.events) > 0 {
		return errpack.New("non flushed events", errpack.Domain())
	}

	e.id = s.Id
	e.createdAt = s.CreatedAt
	e.updatedAt = s.UpdatedAt
	e.deletedAt = s.DeleteAt
	e.version = s.Version

	return nil
}

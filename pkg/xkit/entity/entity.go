package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/xplexer-lab/xplexer/pkg/xkit/binder"
	"github.com/xplexer-lab/xplexer/pkg/xkit/bus"
	"github.com/xplexer-lab/xplexer/pkg/xkit/errpack"
)

var (
	// Base entity has to implement Aggregate
	_ Aggregate[State] = new(Entity)
	_ Entitier         = new(Entity)
	_ Versioner        = new(Entity)
	_ EventProvider    = new(Entity)
	_ Loader[State]    = new(Entity)
	_ Dumper[State]    = new(Entity)
	// ID interfaces
	_ binder.PathUnmarshaler = new(Id)
)

type Id uuid.UUID

func NewId() Id {
	return Id(uuid.New())
}

func ParseId(in string) (Id, error) {
	id, err := uuid.Parse(in)
	return Id(id), err
}

func MustParseId(in string) Id {
	id, err := ParseId(in)

	if err != nil {
		panic(fmt.Errorf("failed to parse id: %w", err))
	}

	return id
}

func (id Id) Hex() string {
	return uuid.UUID(id).String()
}

func (id Id) Validate() error {
	return uuid.Validate(id.Hex())
}

func (id *Id) UnmarshalPath(val string) error {
	if pid, err := ParseId(val); err != nil {
		return err
	} else {
		*id = pid
		return nil
	}
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
	DeletedAt *time.Time
	Version   int
}

type Opt func(*Entity)

func New(opts ...Opt) *Entity {
	return &Entity{
		id:        NewId(),
		createdAt: time.Now().Round(time.Millisecond),
		updatedAt: time.Now().Round(time.Microsecond),
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

func (e *Entity) SetUpdatedAt(t time.Time) {
	e.updatedAt = t
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
		DeletedAt: e.deletedAt,
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
	e.createdAt = s.CreatedAt.Local()
	e.updatedAt = s.UpdatedAt.Local()

	if s.DeletedAt != nil {
		e.deletedAt = lo.ToPtr(s.DeletedAt.Local())
	} else {
		e.deletedAt = nil
	}

	e.version = s.Version

	return nil
}

func (e *Entity) Validate() error {
	// todo: validate other parts
	return e.id.Validate()
}

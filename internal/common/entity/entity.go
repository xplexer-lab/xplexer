package entity

import (
	"time"

	"google.golang.org/protobuf/proto"
)

type EntityId string

type Entity struct {
	id        EntityId
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time

	events []proto.Message
}

type Opt func(*Entity)

func New(opts ...Opt) *Entity {
	return &Entity{
		id:        "",
		createdAt: time.Now(),
		updatedAt: time.Now(),
		deletedAt: nil,
	}
}

func (e *Entity) Id() EntityId {
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

func (e *Entity) RecordEvent(evt proto.Message) {
	e.events = append(e.events, evt)
}

func (e *Entity) PopEvents() []proto.Message {
	ret := e.events
	e.events = nil
	return ret
}

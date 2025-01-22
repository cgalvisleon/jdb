package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
)

type TypeEvent int

const (
	EventInsert TypeEvent = iota
	EventUpdate
	EventDelete
)

func (s TypeEvent) Name() string {
	switch s {
	case EventInsert:
		return "event_insert"
	case EventUpdate:
		return "event_update"
	case EventDelete:
		return "event_delete"
	}
	return ""
}

type Event func(model *Model, before et.Json, after et.Json) error

func EventInsertDefault(model *Model, before et.Json, after et.Json) error {
	event.Publish("model:insert", et.Json{
		"model":  model.Name,
		"table":  model.Table,
		"before": before,
		"after":  after,
	})

	return nil
}

func EventUpdateDefault(model *Model, before et.Json, after et.Json) error {
	event.Publish("model:update", et.Json{
		"model":  model.Name,
		"table":  model.Table,
		"before": before,
		"after":  after,
	})

	return nil
}

func EventDeleteDefault(model *Model, before et.Json, after et.Json) error {
	event.Publish("model:delete", et.Json{
		"model":  model.Name,
		"table":  model.Table,
		"before": before,
		"after":  after,
	})

	return nil
}

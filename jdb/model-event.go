package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/strs"
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

type EventError func(model *Model, data et.Json)

func EventErrorDefault(model *Model, data et.Json) {
	event.Publish("model:error", data)
}

func EventInsertDefault(model *Model, before et.Json, after et.Json) error {
	schema := ""
	if model.Schema != nil {
		schema = model.Schema.Name
	}

	event.Publish("model:insert", et.Json{
		"schema": schema,
		"model":  model.Name,
		"table":  model.Table,
		"before": before,
		"after":  after,
	})

	return nil
}

func EventUpdateDefault(model *Model, before et.Json, after et.Json) error {
	schema := ""
	if model.Schema != nil {
		schema = model.Schema.Name
	}

	event.Publish("model:update", et.Json{
		"schema": schema,
		"model":  model.Name,
		"table":  model.Table,
		"before": before,
		"after":  after,
	})

	return nil
}

func EventDeleteDefault(model *Model, before et.Json, after et.Json) error {
	schema := ""
	if model.Schema != nil {
		schema = model.Schema.Name
	}

	event.Publish("model:delete", et.Json{
		"schema": schema,
		"model":  model.Name,
		"table":  model.Table,
		"before": before,
		"after":  after,
	})

	return nil
}

func EventHistoryDefault(model *Model, before et.Json, after et.Json) error {
	if model.History == nil {
		return nil
	}

	history := model.History
	if history == nil {
		return nil
	}

	if history.With == nil {
		return nil
	}

	key := before.ValStr("", history.Fk.Name)
	if key == "" {
		return nil
	}

	tag := strs.Format("%s:%s", "history", key)
	index := model.Db.GetSerie(tag)
	before[HISTORY_INDEX] = index
	go history.With.Insert(before).
		Exec()

	limit := index - history.Limit
	if limit <= 0 {
		return nil
	}

	go history.With.Delete().
		Where(history.Fk.Name).Eq(key).
		And(HISTORY_INDEX).LessEq(limit).
		Exec()

	return nil
}

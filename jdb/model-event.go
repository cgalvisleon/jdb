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

	if history.Model == nil {
		return nil
	}

	key := before.ValStr("", history.Pk.Name)
	if key == "" {
		return nil
	}

	key = strs.Format("%s:%s", "history", key)
	index := model.Db.GetSerie(key)
	before[HISTORY_INDEX] = index
	go history.Model.Insert(before).
		Exec()

	limit := index - model.History.Limit
	if limit <= 0 {
		return nil
	}

	go history.Model.Delete().
		Where(history.Pk.Name).Eq(key).
		And(HISTORY_INDEX).LessEq(limit).
		Exec()

	return nil
}

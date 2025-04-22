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

/**
* EventErrorDefault
* @param model *Model, err et.Json
**/
func EventErrorDefault(model *Model, err et.Json) {
	schema := ""
	if model.Schema != nil {
		schema = model.Schema.Name
	}

	data := et.Json{
		"schema": schema,
		"model":  model.Name,
		"table":  model.Table,
		"error":  err,
	}

	event.Publish(event.EVENT_MODEL_ERROR, data)
	event.Publish(event.EVENT_MODEL_ERROR+model.Name, data)
	event.Publish(event.EVENT_MODEL_ERROR+model.Table, data)
}

/**
* EventInsertDefault
* @param model *Model, before et.Json, after et.Json
* @return error
**/
func EventInsertDefault(model *Model, before et.Json, after et.Json) error {
	schema := ""
	if model.Schema != nil {
		schema = model.Schema.Name
	}

	data := et.Json{
		"schema": schema,
		"model":  model.Name,
		"table":  model.Table,
		"before": before,
		"after":  after,
	}

	event.Publish(event.EVENT_MODEL_INSERT, data)
	event.Publish(event.EVENT_MODEL_INSERT+model.Name, data)
	event.Publish(event.EVENT_MODEL_INSERT+model.Table, data)

	return nil
}

/**
* EventUpdateDefault
* @param model *Model, before et.Json, after et.Json
* @return error
**/
func EventUpdateDefault(model *Model, before et.Json, after et.Json) error {
	schema := ""
	if model.Schema != nil {
		schema = model.Schema.Name
	}

	if model.StatusField != nil && model.SystemKeyField != nil {
		oldStatus := before.Str(model.StatusField.Name)
		newStatus := after.Str(model.StatusField.Name)
		if oldStatus != newStatus {
			model.Db.upsertRecycling(model.Table, before.Str(model.SystemKeyField.Name), newStatus)
		}
	}

	data := et.Json{
		"schema": schema,
		"model":  model.Name,
		"table":  model.Table,
		"before": before,
		"after":  after,
	}

	event.Publish(event.EVENT_MODEL_UPDATE, data)
	event.Publish(event.EVENT_MODEL_UPDATE+model.Name, data)
	event.Publish(event.EVENT_MODEL_UPDATE+model.Table, data)

	return nil
}

/**
* EventDeleteDefault
* @param model *Model, before et.Json, after et.Json
* @return error
**/
func EventDeleteDefault(model *Model, before et.Json, after et.Json) error {
	schema := ""
	if model.Schema != nil {
		schema = model.Schema.Name
	}

	data := et.Json{
		"schema": schema,
		"model":  model.Name,
		"table":  model.Table,
		"before": before,
		"after":  after,
	}

	event.Publish(event.EVENT_MODEL_DELETE, data)
	event.Publish(event.EVENT_MODEL_DELETE+model.Name, data)
	event.Publish(event.EVENT_MODEL_DELETE+model.Table, data)

	return nil
}

/**
* EventHistoryDefault
* @param model *Model, before et.Json, after et.Json
* @return error
**/
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

	limit := index - int64(history.Limit)
	if limit <= 0 {
		return nil
	}

	go history.With.Delete().
		Where(history.Fk.Name).Eq(key).
		And(HISTORY_INDEX).LessEq(limit).
		Exec()

	return nil
}

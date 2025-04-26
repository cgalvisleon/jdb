package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/strs"
)

const EVENT_MODEL_ERROR = "model:error"
const EVENT_MODEL_INSERT = "model:insert"
const EVENT_MODEL_UPDATE = "model:update"
const EVENT_MODEL_DELETE = "model:delete"

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
* eventErrorDefault
* @param model *Model, err et.Json
**/
func eventErrorDefault(model *Model, err et.Json) {
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

	event.Publish(EVENT_MODEL_ERROR, data)
	event.Publish(EVENT_MODEL_ERROR+model.Name, data)
	event.Publish(EVENT_MODEL_ERROR+model.Table, data)
}

/**
* eventInsertDefault
* @param model *Model, before et.Json, after et.Json
* @return error
**/
func eventInsertDefault(model *Model, before et.Json, after et.Json) error {
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

	event.Publish(EVENT_MODEL_INSERT, data)
	event.Publish(EVENT_MODEL_INSERT+model.Name, data)
	event.Publish(EVENT_MODEL_INSERT+model.Table, data)

	return nil
}

/**
* eventUpdateDefault
* @param model *Model, before et.Json, after et.Json
* @return error
**/
func eventUpdateDefault(model *Model, before et.Json, after et.Json) error {
	schema := ""
	if model.Schema != nil {
		schema = model.Schema.Name
	}

	if model.StatusField != nil && model.SystemKeyField != nil {
		oldStatus := before.Str(model.StatusField.Name)
		newStatus := after.Str(model.StatusField.Name)
		if oldStatus != newStatus {
			sysId := before.Str(model.SystemKeyField.Name)
			model.Db.upsertRecycling(model.Table, sysId, newStatus)
		}
	}

	data := et.Json{
		"schema": schema,
		"model":  model.Name,
		"table":  model.Table,
		"before": before,
		"after":  after,
	}

	event.Publish(EVENT_MODEL_UPDATE, data)
	event.Publish(EVENT_MODEL_UPDATE+model.Name, data)
	event.Publish(EVENT_MODEL_UPDATE+model.Table, data)

	return nil
}

/**
* eventDeleteDefault
* @param model *Model, before et.Json, after et.Json
* @return error
**/
func eventDeleteDefault(model *Model, before et.Json, after et.Json) error {
	schema := ""
	if model.Schema != nil {
		schema = model.Schema.Name
	}

	if model.StatusField != nil && model.SystemKeyField != nil {
		sysId := before.Str(model.SystemKeyField.Name)
		model.Db.deleteRecycling(model.Table, sysId)
	}

	data := et.Json{
		"schema": schema,
		"model":  model.Name,
		"table":  model.Table,
		"before": before,
		"after":  after,
	}

	event.Publish(EVENT_MODEL_DELETE, data)
	event.Publish(EVENT_MODEL_DELETE+model.Name, data)
	event.Publish(EVENT_MODEL_DELETE+model.Table, data)

	return nil
}

/**
* eventHistoryDefault
* @param model *Model, before et.Json
* @return error
**/
func eventHistoryDefault(model *Model, before et.Json) error {
	if model.History == nil {
		return nil
	}

	history := model.History.With
	if history == nil {
		return nil
	}

	tag := "history"
	n := 0
	command := history.
		Insert(before)
	for fkn, pk := range model.History.Fk {
		key := before.ValStr("", pk)
		if n == 0 {
			command.Where(fkn).Eq(key)
		} else {
			command.And(fkn).Eq(key)
		}
		tag = strs.Append(tag, key, ":")
		n++
	}

	index, err := model.Db.GetSerie(tag)
	if err != nil {
		return err
	}
	before[HISTORY_INDEX] = index
	go command.Exec()

	limit := index - int64(model.History.Limit)
	if limit <= 0 {
		return nil
	}

	n = 0
	command = history.
		Delete()
	for fkn, pk := range model.History.Fk {
		key := before.ValStr("", pk)
		if n == 0 {
			command.Where(fkn).Eq(key)
		} else {
			command.And(fkn).Eq(key)
		}
		n++
	}
	command.And(HISTORY_INDEX).LessEq(limit)
	go command.Exec()

	return nil
}

package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/utility"
)

const EVENT_MODEL_ERROR = "model:error"
const EVENT_MODEL_INSERT = "model:insert"
const EVENT_MODEL_UPDATE = "model:update"
const EVENT_MODEL_DELETE = "model:delete"

/**
* publishError
* @param model *Model, sql string, err error
**/
func publishError(model *Model, sql string, err error) {
	event.Publish(EVENT_MODEL_ERROR, et.Json{
		"db":     model.Db.Name,
		"schema": model.Schema,
		"model":  model.Name,
		"sql":    sql,
		"error":  err.Error(),
	})
}

/**
* publishInsert
* @param model *Model, sql string
**/
func publishInsert(model *Model, sql string) {
	event.Publish(EVENT_MODEL_INSERT, et.Json{
		"db":      model.Db.Name,
		"schema":  model.Schema,
		"model":   model.Name,
		"sql":     sql,
		"command": "insert",
	})
}

/**
* publishUpdate
* @param model *Model, sql string
**/
func publishUpdate(model *Model, sql string) {
	event.Publish(EVENT_MODEL_UPDATE, et.Json{
		"db":      model.Db.Name,
		"schema":  model.Schema,
		"model":   model.Name,
		"sql":     sql,
		"command": "update",
	})
}

/**
* publishDelete
* @param model *Model, sql string
**/
func publishDelete(model *Model, sql string) {
	event.Publish(EVENT_MODEL_DELETE, et.Json{
		"db":      model.Db.Name,
		"schema":  model.Schema,
		"model":   model.Name,
		"sql":     sql,
		"command": "delete",
	})
}

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

type Event func(tx *Tx, model *Model, before et.Json, after et.Json) error

/**
* eventInsertDefault
* @param tx *Tx, model *Model, before et.Json, after et.Json
* @return error
**/
func eventInsertDefault(tx *Tx, model *Model, before et.Json, after et.Json) error {
	if model.UseCore {
		model.Db.upsertTable(tx, model.Schema, model.Name, 1)

		if model.SystemKeyField != nil {
			sysId := after.Str(model.SystemKeyField.Name)
			err := model.Db.upsertRecord(tx, model.Schema, model.Name, sysId, "insert")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

/**
* eventUpdateDefault
* @param tx *Tx, model *Model, before et.Json, after et.Json
* @return error
**/
func eventUpdateDefault(tx *Tx, model *Model, before et.Json, after et.Json) error {
	if model.UseCore && model.SystemKeyField != nil {
		sysId := after.Str(model.SystemKeyField.Name)
		err := model.Db.upsertRecord(tx, model.Schema, model.Name, sysId, "update")
		if err != nil {
			return err
		}

		if model.StatusField != nil {
			oldStatus := before.ValStr(utility.ACTIVE, model.StatusField.Name)
			newStatus := after.ValStr(oldStatus, model.StatusField.Name)
			if oldStatus != newStatus && newStatus == utility.FOR_DELETE {
				err := model.Db.upsertRecycling(tx, model.Schema, model.Name, sysId)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

/**
* eventDeleteDefault
* @param tx *Tx, model *Model, before et.Json, after et.Json
* @return error
**/
func eventDeleteDefault(tx *Tx, model *Model, before et.Json, after et.Json) error {
	if model.UseCore {
		model.Db.upsertTable(tx, model.Schema, model.Name, -1)

		if model.SystemKeyField != nil {
			sysId := before.Str(model.SystemKeyField.Name)
			err := model.Db.upsertRecord(tx, model.Schema, model.Name, sysId, "delete")
			if err != nil {
				return err
			}

			if model.StatusField != nil {
				err := model.Db.deleteRecycling(tx, model.Schema, model.Name, sysId)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

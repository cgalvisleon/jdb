package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

const EVENT_MODEL_ERROR = "model:error"
const EVENT_MODEL_INSERT = "model:insert"
const EVENT_MODEL_UPDATE = "model:update"
const EVENT_MODEL_DELETE = "model:delete"
const EVENT_MODEL_SYNC = "model:sync"

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

type EventError func(model *Model, data et.Json)

/**
* eventErrorDefault
* @param model *Model, err et.Json
**/
func eventErrorDefault(model *Model, err et.Json) {
	data := et.Json{
		"db":     model.Db.Name,
		"schema": model.Schema,
		"model":  model.Name,
		"error":  err,
	}

	event.Publish(fmt.Sprintf("%s:%s", EVENT_MODEL_ERROR, model.Db.Name), data)
}

/**
* eventInsertDefault
* @param tx *Tx, model *Model, before et.Json, after et.Json
* @return error
**/
func eventInsertDefault(tx *Tx, model *Model, before et.Json, after et.Json) error {
	if model.UseCore && model.SystemKeyField != nil {
		sysid := after.Str(model.SystemKeyField.Name)
		err := model.Db.upsertRecord(tx, model.Schema, model.Name, sysid, "insert")
		if err != nil {
			return err
		}
	}

	data := et.Json{
		"db":     model.Db.Name,
		"schema": model.Schema,
		"model":  model.Name,
		"before": before,
		"after":  after,
	}

	event.Publish(fmt.Sprintf("%s:%s", EVENT_MODEL_INSERT, model.Db.Name), data)

	return nil
}

/**
* eventUpdateDefault
* @param tx *Tx, model *Model, before et.Json, after et.Json
* @return error
**/
func eventUpdateDefault(tx *Tx, model *Model, before et.Json, after et.Json) error {
	if model.UseCore && model.SystemKeyField != nil && !model.isAudit {
		sysid := after.Str(model.SystemKeyField.Name)
		err := model.Db.upsertRecord(tx, model.Schema, model.Name, sysid, "update")
		if err != nil {
			return err
		}
	}

	if model.UseCore && model.SystemKeyField != nil && model.StatusField != nil {
		oldStatus := before.ValStr(utility.ACTIVE, model.StatusField.Name)
		newStatus := after.ValStr(oldStatus, model.StatusField.Name)
		if oldStatus != newStatus {
			sysId := before.Str(model.SystemKeyField.Name)
			err := model.Db.upsertRecycling(tx, model.Schema, model.Name, sysId, newStatus)
			if err != nil {
				return err
			}
		}
	}

	data := et.Json{
		"db":     model.Db.Name,
		"schema": model.Schema,
		"model":  model.Name,
		"before": before,
		"after":  after,
	}

	event.Publish(fmt.Sprintf("%s:%s", EVENT_MODEL_UPDATE, model.Db.Name), data)

	return nil
}

/**
* eventDeleteDefault
* @param tx *Tx, model *Model, before et.Json, after et.Json
* @return error
**/
func eventDeleteDefault(tx *Tx, model *Model, before et.Json, after et.Json) error {
	if model.UseCore && model.SystemKeyField != nil && !model.isAudit {
		sysid := after.Str(model.SystemKeyField.Name)
		err := model.Db.upsertRecord(tx, model.Schema, model.Name, sysid, "delete")
		if err != nil {
			return err
		}
	}

	if model.UseCore && model.SystemKeyField != nil && model.StatusField != nil {
		sysId := before.Str(model.SystemKeyField.Name)
		err := model.Db.deleteRecycling(tx, model.Schema, model.Name, sysId)
		if err != nil {
			return err
		}
	}

	data := et.Json{
		"db":     model.Db.Name,
		"schema": model.Schema,
		"model":  model.Name,
		"before": before,
		"after":  after,
	}

	event.Publish(fmt.Sprintf("%s:%s", EVENT_MODEL_DELETE, model.Db.Name), data)

	return nil
}

/**
* eventSyncDefault
* @param message event.Message
**/
func eventSyncDefault(message event.Message) {
	data := message.Data
	db := data.Str("db")
	schema := data.Str("schema")
	model := data.Str("model")
	syncChannel := strs.Format("sync:%s.%s.%s", db, schema, model)
	event.Publish(syncChannel, data)
}

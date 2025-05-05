package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/timezone"
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
* eventHistoryDefault
* @param tx *Tx, model *Model, before et.Json
* @return error
**/
func eventHistoryDefault(tx *Tx, model *Model, before et.Json) error {
	if model.History == nil {
		return nil
	}

	if model.SystemKeyField == nil {
		return nil
	}

	history := model.History.With
	if history == nil {
		return nil
	}

	sysId := before.Str(model.SystemKeyField.Name)
	if sysId == "" {
		return nil
	}

	index, err := model.Db.GetSerie(sysId)
	if err != nil {
		return err
	}

	_, err = history.
		Insert(et.Json{
			CREATED_AT: timezone.Now(),
			SYSID:      sysId,
			HISTORYCAL: before,
			INDEX:      index,
		}).
		ExecTx(tx)
	if err != nil {
		return err
	}

	limit := index - int64(model.History.Limit)
	if limit <= 0 {
		return nil
	}

	_, err = history.
		Delete(SYSID).Eq(sysId).
		And(INDEX).LessEq(limit).
		ExecTx(tx)
	if err != nil {
		return err
	}

	return nil
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
	if model.UseCore && model.SystemKeyField != nil {
		sysid := after.Str(model.SystemKeyField.Name)
		err := model.Db.upsertRecord(tx, model.Schema, model.Name, sysid, "update")
		if err != nil {
			return err
		}
	}

	if model.UseCore && model.SystemKeyField != nil && model.StatusField != nil {
		oldStatus := before.Str(model.StatusField.Name)
		newStatus := after.Str(model.StatusField.Name)
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
	if model.UseCore && model.SystemKeyField != nil {
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

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

func EventHistoryDefault(model *Model, before et.Json, after et.Json) error {
	if model.History == nil {
		return nil
	}

	key := before.Key(model.KeyField.Name)
	if key == "-1" {
		return nil
	}

	index := model.Db.GetSerie(key)
	before[HISTORY_INDEX] = index
	go model.History.Insert(before).
		Exec()

	limit := index - model.HistoryLimit
	if limit <= 0 {
		return nil
	}

	go model.History.Delete().
		Where(model.KeyField.Name).Eq(key).
		And(HISTORY_INDEX).LessEq(limit).
		Exec()

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

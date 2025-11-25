package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
)

func (s *Model) publish(action string, data et.Json) {
	key := s.GetKey(data)
	recordId := data.String(s.RecordField)
	if key != "" {
		event.Publish(key, et.Json{
			"database": s.Database,
			"model":    s.Name,
			"action":   action,
			KEY:        key,
			RECORDID:   recordId,
			"data":     data,
		})
	}

	if recordId != "" {
		event.Publish(key, et.Json{
			"database": s.Database,
			"model":    s.Name,
			"action":   action,
			KEY:        key,
			RECORDID:   recordId,
			"data":     data,
		})
	}

	channel := fmt.Sprintf("%s:%s", s.Database, s.Name)
	if channel != "" {
		event.Publish(channel, et.Json{
			"database": s.Database,
			"model":    s.Name,
			"action":   action,
			KEY:        key,
			RECORDID:   recordId,
			"data":     data,
		})
	}
}

/**
* afterInsertDefault
* @param tx *Tx, old, new et.Json
* @return error
**/
func (s *Model) afterInsertDefault(tx *Tx, old, new et.Json) error {
	s.publish("insert", new)

	if s.isCore {
		return nil
	}

	return nil
}

/**
* afterUpdateDefault
* @param tx *Tx, data et.Json
* @return error
**/
func (s *Model) afterUpdateDefault(tx *Tx, old, new et.Json) error {
	s.publish("update", new)

	if s.isCore {
		return nil
	}

	return nil
}

/**
* afterDeleteDefault
* @param tx *Tx, data et.Json
* @return error
**/
func (s *Model) afterDeleteDefault(tx *Tx, old, new et.Json) error {
	s.publish("delete", new)

	if s.isCore {
		return nil
	}

	return nil
}

/**
* AfterInsert
* @param fn DataFunction
* @return *Command
**/
func (s *Model) AfterInsert(fn DataFunctionTx) *Model {
	s.afterInserts = append(s.afterInserts, fn)

	return s
}

/**
* AfterUpdate
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Model) AfterUpdate(fn DataFunctionTx) *Model {
	s.afterUpdates = append(s.afterUpdates, fn)

	return s
}

/**
* AfterDelete
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Model) AfterDelete(fn DataFunctionTx) *Model {
	s.afterDeletes = append(s.afterDeletes, fn)

	return s
}

/**
* AfterInsertOrUpdate
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Model) AfterInsertOrUpdate(fn DataFunctionTx) *Model {
	s.afterInserts = append(s.afterInserts, fn)
	s.afterUpdates = append(s.afterUpdates, fn)

	return s
}

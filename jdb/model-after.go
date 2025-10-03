package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
)

/**
* afterInsertDefault
* @param tx *Tx, data et.Json
* @return error
**/
func (s *Model) afterInsertDefault(tx *Tx, data et.Json) error {
	if s.RecordField != "" {
		id := data.String(s.RecordField)
		event.Publish(id, et.Json{
			"database": s.Database,
			"model":    s.Name,
			"action":   "insert",
			RECORDID:   id,
			"data":     data,
		})
		key := fmt.Sprintf("%s:%s", s.Database, s.Name)
		event.Publish(key, et.Json{
			"database": s.Database,
			"model":    s.Name,
			"action":   "insert",
			RECORDID:   id,
			"data":     data,
		})
	}

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
func (s *Model) afterUpdateDefault(tx *Tx, data et.Json) error {
	if s.RecordField != "" {
		id := data.String(s.RecordField)
		event.Publish(id, et.Json{
			"database": s.Database,
			"model":    s.Name,
			"action":   "update",
			RECORDID:   id,
			"data":     data,
		})
		key := fmt.Sprintf("%s:%s", s.Database, s.Name)
		event.Publish(key, et.Json{
			"database": s.Database,
			"model":    s.Name,
			"action":   "update",
			RECORDID:   id,
			"data":     data,
		})
	}

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
func (s *Model) afterDeleteDefault(tx *Tx, data et.Json) error {
	if s.RecordField != "" {
		id := data.String(s.RecordField)
		event.Publish(id, et.Json{
			"database": s.Database,
			"model":    s.Name,
			"action":   "delete",
			RECORDID:   id,
			"data":     data,
		})
		key := fmt.Sprintf("%s:%s", s.Database, s.Name)
		event.Publish(key, et.Json{
			"database": s.Database,
			"model":    s.Name,
			"action":   "delete",
			RECORDID:   id,
			"data":     data,
		})
	}

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

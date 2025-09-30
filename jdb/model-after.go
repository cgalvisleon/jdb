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
func (s *Model) EventAfterInsert(fn DataFunctionTx) *Model {
	s.eventAfterInsert = append(s.eventAfterInsert, fn)

	return s
}

/**
* EventAfterUpdate
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Model) EventAfterUpdate(fn DataFunctionTx) *Model {
	s.eventAfterUpdate = append(s.eventAfterUpdate, fn)

	return s
}

/**
* EventAfterDelete
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Model) EventAfterDelete(fn DataFunctionTx) *Model {
	s.eventAfterDelete = append(s.eventAfterDelete, fn)

	return s
}

/**
* EventAfterInsertOrUpdate
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Model) EventAfterInsertOrUpdate(fn DataFunctionTx) *Model {
	s.eventAfterInsert = append(s.eventAfterInsert, fn)
	s.eventAfterUpdate = append(s.eventAfterUpdate, fn)

	return s
}

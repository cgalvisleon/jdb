package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/reg"
)

/**
* beforeInsertDefault
* @param tx *Tx, data et.Json
* @return error
**/
func (s *Model) beforeInsertDefault(tx *Tx, data et.Json) error {
	if s.RecordField != "" {
		id := reg.GenULIDI(s.Name)
		data.Set(s.RecordField, id)
	}

	if s.isCore {
		return nil
	}

	return nil
}

/**
* beforeUpdateDefault
* @param tx *Tx, data et.Json
* @return error
**/
func (s *Model) beforeUpdateDefault(tx *Tx, data et.Json) error {
	if s.isCore {
		return nil
	}

	return nil
}

/**
* beforeDeleteDefault
* @param tx *Tx, data et.Json
* @return error
**/
func (s *Model) beforeDeleteDefault(tx *Tx, data et.Json) error {
	if s.isCore {
		return nil
	}

	return nil
}

/**
* EventBeforeInsert
* @param fn DataFunction
**/
func (s *Model) EventBeforeInsert(fn DataFunctionTx) *Model {
	s.eventBeforeInsert = append(s.eventBeforeInsert, fn)

	return s
}

/**
* EventBeforeUpdate
* @param fn DataFunction
**/
func (s *Model) EventBeforeUpdate(fn DataFunctionTx) *Model {
	s.eventBeforeUpdate = append(s.eventBeforeUpdate, fn)

	return s
}

/**
* EventBeforeDelete
* @param fn DataFunction
**/
func (s *Model) EventBeforeDelete(fn DataFunctionTx) *Model {
	s.eventBeforeDelete = append(s.eventBeforeDelete, fn)

	return s
}

/**
* EventBeforeInsertOrUpdate
* @param fn DataFunction
**/
func (s *Model) EventBeforeInsertOrUpdate(fn DataFunctionTx) *Model {
	s.eventBeforeInsert = append(s.eventBeforeInsert, fn)
	s.eventBeforeUpdate = append(s.eventBeforeUpdate, fn)

	return s
}

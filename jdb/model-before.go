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
* BeforeInsert
* @param fn DataFunction
**/
func (s *Model) BeforeInsert(fn DataFunctionTx) *Model {
	s.beforeInserts = append(s.beforeInserts, fn)

	return s
}

/**
* BeforeUpdate
* @param fn DataFunction
**/
func (s *Model) BeforeUpdate(fn DataFunctionTx) *Model {
	s.beforeUpdates = append(s.beforeUpdates, fn)

	return s
}

/**
* BeforeDelete
* @param fn DataFunction
**/
func (s *Model) BeforeDelete(fn DataFunctionTx) *Model {
	s.beforeDeletes = append(s.beforeDeletes, fn)

	return s
}

/**
* BeforeInsertOrUpdate
* @param fn DataFunction
**/
func (s *Model) BeforeInsertOrUpdate(fn DataFunctionTx) *Model {
	s.beforeInserts = append(s.beforeInserts, fn)
	s.beforeUpdates = append(s.beforeUpdates, fn)

	return s
}

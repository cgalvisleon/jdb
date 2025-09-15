package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/utility"
)

/**
* beforeInsertDefault
* @param tx *Tx, data et.Json
* @return error
**/
func (s *Command) beforeInsertDefault(tx *Tx, data et.Json) error {
	if s.From == nil {
		return fmt.Errorf(MSG_MODEL_REQUIRED)
	}

	model := s.From

	if model.IndexField != nil && data.Int(model.IndexField.Name) == 0 {
		data[model.IndexField.Name] = reg.GenIndex()
	}

	if model.SystemKeyField != nil && data.Str(model.SystemKeyField.Name) == "" {
		data[model.SystemKeyField.Name] = model.GenId()
	}

	now := utility.Now()
	if model.CreatedAtField != nil && data.Str(model.CreatedAtField.Name) == "" {
		data[model.CreatedAtField.Name] = now
	}

	if model.UpdatedAtField != nil && data.Str(model.UpdatedAtField.Name) == "" {
		data[model.UpdatedAtField.Name] = now
	}

	return nil
}

/**
* beforeUpdateDefault
* @param tx *Tx, data et.Json
* @return error
**/
func (s *Command) beforeUpdateDefault(tx *Tx, data et.Json) error {
	if s.From == nil {
		return fmt.Errorf(MSG_MODEL_REQUIRED)
	}

	now := utility.Now()
	model := s.From
	if model.CreatedAtField != nil {
		delete(data, model.CreatedAtField.Name)
	}

	if model.UpdatedAtField != nil && data.Str(model.UpdatedAtField.Name) == "" {
		data[model.UpdatedAtField.Name] = now
	}

	return nil
}

/**
* beforeDeleteDefault
* @param tx *Tx, data et.Json
* @return error
**/
func (s *Command) beforeDeleteDefault(tx *Tx, data et.Json) error {
	if s.From == nil {
		return fmt.Errorf(MSG_MODEL_REQUIRED)
	}

	return nil
}

/**
* BeforeInsert
* @param fn DataFunction
**/
func (s *Command) BeforeInsert(fn DataFunctionTx) *Command {
	s.beforeInsert = append(s.beforeInsert, fn)

	return s
}

/**
* BeforeUpdate
* @param fn DataFunction
**/
func (s *Command) BeforeUpdate(fn DataFunctionTx) *Command {
	s.beforeUpdate = append(s.beforeUpdate, fn)

	return s
}

/**
* BeforeDelete
* @param fn DataFunction
**/
func (s *Command) BeforeDelete(fn DataFunctionTx) *Command {
	s.beforeDelete = append(s.beforeDelete, fn)

	return s
}

/**
* BeforeInsertOrUpdate
* @param fn DataFunction
**/
func (s *Command) BeforeInsertOrUpdate(fn DataFunctionTx) *Command {
	s.beforeInsert = append(s.beforeInsert, fn)
	s.beforeUpdate = append(s.beforeUpdate, fn)

	return s
}

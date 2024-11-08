package jdb

import (
	"slices"
)

/**
* DefineColumn
* @param name string
* @param typeData TypeData
* @param def interface{}
* @return *Model
**/
func (s *Model) DefineColumn(name string, typeData TypeData, def interface{}) *Column {
	col := newColumn(s, name, "", TpColumn, typeData, def)
	s.Columns[col.Field] = col
	if slices.Contains([]string{IndexField, ProjectField, CreatedAtField, UpdatedAtField, StateField, SystemKeyField, FullTextField}, name) {
		s.DefineIndex(true, name)
	}

	return col
}

/**
* DefineAtribute
* @param name string
* @param typeData TypeData
* @param def interface{}
* @return *Model
**/
func (s *Model) DefineAtribute(name string, typeData TypeData, def interface{}) *Column {
	col := newColumn(s, name, "", TpAtribute, typeData, def)
	s.Columns[col.Field] = col

	return col
}

/**
* DefineKey
* @param colums ...*Column
* @return *Model
**/
func (s *Model) DefineKey(colums ...string) *Model {
	cols := s.GetColumns(colums...)
	if len(cols) == 0 {
		return s
	}

	for _, col := range cols {
		s.Keys[col.Field] = col
	}

	return s
}

/**
* DefineIndex
* @param colums ...*Column
* @return *Model
**/
func (s *Model) DefineIndex(sort bool, colums ...string) *Model {
	cols := s.GetColumns(colums...)
	if len(cols) == 0 {
		return s
	}

	for _, col := range cols {
		idx := NewIndex(col, sort)
		s.Indices[col.Field] = idx
	}

	return s
}

/**
* DefineUnique
* @param colums ...*Column
* @return *Model
**/
func (s *Model) DefineUnique(colums ...string) *Model {
	cols := s.GetColumns(colums...)
	if len(cols) == 0 {
		return s
	}

	for _, col := range cols {
		idx := NewIndex(col, true)
		s.Uniques[col.Field] = idx
	}

	return s
}

/**
* DefineTrigger
* @param tp TypeTrigger
* @param trigger Trigger
* @return *Model
**/
func (s *Model) DefineTrigger(tp TypeTrigger, trigger Trigger) *Model {
	switch tp {
	case BeforeInsert:
		s.BeforeInsert = append(s.BeforeInsert, trigger)
	case AfterInsert:
		s.AfterInsert = append(s.AfterInsert, trigger)
	case BeforeUpdate:
		s.BeforeUpdate = append(s.BeforeUpdate, trigger)
	case AfterUpdate:
		s.AfterUpdate = append(s.AfterUpdate, trigger)
	case BeforeDelete:
		s.BeforeDelete = append(s.BeforeDelete, trigger)
	case AfterDelete:
		s.AfterDelete = append(s.AfterDelete, trigger)
	}
	return s
}

/**
* DefineFunction
* @param name string
* @param function string
* @return *Model
**/
func (s *Model) DefineFunction(name string, tp TypeFunction, definition string) *Function {
	f := NewFunction(name, tp)
	f.Definition = definition
	s.Functions[f.Key] = f

	return f
}

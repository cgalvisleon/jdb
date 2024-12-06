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
func (s *Model) DefineColumn(name string, typeData TypeData) *Column {
	def := typeData.DefaultValue(s.Db.driver)
	col := newColumn(s, name, "", TpColumn, typeData, def)
	s.Columns = append(s.Columns, col)
	if slices.Contains([]string{IndexField.Str(), ProjectField.Str(), CreatedAtField.Str(), UpdatedAtField.Str(), StateField.Str(), KeyField.Str(), SystemKeyField.Str(), DataField.Str(), FullTextField.Str()}, name) {
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
func (s *Model) DefineAtribute(name string, typeData TypeData) *Column {
	def := typeData.DefaultValue(s.Db.driver)
	col := newColumn(s, name, "", TpAtribute, typeData, def)
	s.Columns = append(s.Columns, col)

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
* DefineRequired
* @param requireds ...string
* @return *Model
**/
func (s *Model) DefineRequired(requireds ...string) *Model {
	for _, required := range requireds {
		col := s.GetColumn(required)
		if col != nil {
			s.ColRequired[col.Field] = true
		}
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
func (s *Model) DefineFunction(name string, tp TypeFunction, definition string) *Model {
	f := NewFunction(name, tp)
	f.Definition = definition
	s.Functions[f.Key] = f

	return s
}

/**
* DefineDetail
* @param name string
* @param detail Detail
* @return *Model
**/
func (s *Model) DefineDetail(name string, detail Detail) *Model {
	s.Details[name] = detail

	return s
}

/**
* DefineDictionary
* @param name string
* @param key string
* @param value interface{}
* @return *Dictionary
**/
func (s *Model) DefineDictionary(name, key string, value interface{}) *Dictionary {
	result := NewDictionary(s, key, value)
	s.Dictionaries[name] = result

	return result
}

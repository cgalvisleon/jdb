package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/strs"
)

/**
* DefineColumn
* @param name string
* @param typeData TypeData
* @param def interface{}
* @return *Model
**/
func (s *Model) DefineColumn(name string, typeData TypeData) *Column {
	def := typeData.DefaultValue()
	col := newColumn(s, name, "", TpColumn, typeData, def)
	s.Columns = append(s.Columns, col)
	if strs.Uppcase(col.Name) == SourceField.Uppcase() {
		s.SourceField = col
	}
	if strs.Uppcase(col.Name) == SystemKeyField.Uppcase() {
		s.SystemKeyField = col
	}
	if strs.Uppcase(col.Name) == IndexField.Uppcase() {
		s.IndexField = col
	}
	if strs.Uppcase(col.Name) == StateField.Uppcase() {
		s.StateField = col
	}
	if strs.Uppcase(col.Name) == ClassField.Uppcase() {
		s.ClassField = col
	}
	if slices.Contains([]string{IndexField.Str(), ProjectField.Str(), CreatedAtField.Str(), UpdatedAtField.Str(), StateField.Str(), KeyField.Str(), SystemKeyField.Str(), SourceField.Str(), FullTextField.Str()}, name) {
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
	def := typeData.DefaultValue()
	col := newColumn(s, name, "", TpAtribute, typeData, def)
	s.Columns = append(s.Columns, col)

	return col
}

/**
* DefineGenerate
* @param name string
* @param function string
* @return *Model
**/
func (s *Model) DefineGenerated(name string, f FuncGenerated) {
	col := newColumn(s, name, "", TpGenerate, TypeDataNone, TypeDataNone.DefaultValue())
	col.Definition = f
	s.Columns = append(s.Columns, col)
}

/**
* DefineDetail
* @param name string
* @param detail Detail
* @return *Model
**/
func (s *Model) DefineDetail(name string) *Model {
	detail := NewModel(s.Schema, name)
	col := newColumn(s, name, "", TpDetail, TypeDataNone, TypeDataNone.DefaultValue())
	col.Definition = detail
	s.Columns = append(s.Columns, col)
	keys := s.GetKeys()
	for _, key := range keys {
		fkn := key.Fk()
		fk := detail.DefineColumn(fkn, key.TypeData)
		NewReference(fk, RelationManyToOne, key)
		ref := NewReference(key, RelationOneToMany, fk)
		ref.OnDeleteCascade = true
		ref.OnUpdateCascade = true
	}

	return detail
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
func (s *Model) DefineFunction(name string, tp TypeFunction, definition string) *Function {
	f := NewFunction(name, tp)
	f.Definition = definition
	s.Functions[f.Id] = f

	return f
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

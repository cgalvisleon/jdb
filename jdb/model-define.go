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
	col := s.GetColumn(name)
	if col != nil {
		return col
	}

	def := typeData.DefaultValue()
	col = newColumn(s, name, "", TpColumn, typeData, def)
	if col.Up() == SourceField.Up() {
		s.SourceField = col
	}
	if col.Up() == SystemKeyField.Up() {
		s.SystemKeyField = col
	}
	if col.Up() == IndexField.Up() {
		s.IndexField = col
	}
	if col.Up() == StateField.Up() {
		s.StateField = col
	}
	if col.Up() == ClassField.Up() {
		s.ClassField = col
	}
	if col.Up() == IndexField.Up() {
		s.IndexField = col
	}
	idx := slices.IndexFunc(s.Columns, func(e *Column) bool { return e == s.IndexField })
	if idx == -1 {
		s.Columns = append(s.Columns, col)
	} else {
		s.Columns = append(s.Columns[:idx], append([]*Column{col}, s.Columns[idx:]...)...)
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
	col := s.GetColumn(name)
	if col != nil {
		return col
	}

	s.DefineColumn(SourceField.Up(), SourceField.TypeData())
	def := typeData.DefaultValue()
	col = newColumn(s, name, "", TpAtribute, typeData, def)
	col.Field = SourceField.Low()
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
func (s *Model) DefineDetail(name, table string, version int) *Model {
	detail := NewModel(s.Schema, table, version)
	col := newColumn(s, name, "", TpDetail, TypeDataNone, TypeDataNone.DefaultValue())
	col.Definition = detail
	s.Columns = append(s.Columns, col)
	keys := s.GetKeys()
	detail.MakeDetail(keys)

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
	name := strs.Format(`%s_%s`, s.Name, tp.Name())
	name = strs.Uppcase(name)
	switch tp {
	case BeforeInsert:
		idx := len(s.BeforeInsert) + 1
		name := strs.Format(`%s_%d`, name, idx)
		Triggers[name] = trigger
		s.BeforeInsert = append(s.BeforeInsert, name)
	case AfterInsert:
		idx := len(s.AfterInsert) + 1
		name := strs.Format(`%s_%d`, name, idx)
		Triggers[name] = trigger
		s.AfterInsert = append(s.AfterInsert, name)
	case BeforeUpdate:
		idx := len(s.BeforeUpdate) + 1
		name := strs.Format(`%s_%d`, name, idx)
		Triggers[name] = trigger
		s.BeforeUpdate = append(s.BeforeUpdate, name)
	case AfterUpdate:
		idx := len(s.AfterUpdate) + 1
		name := strs.Format(`%s_%d`, name, idx)
		Triggers[name] = trigger
		s.AfterUpdate = append(s.AfterUpdate, name)
	case BeforeDelete:
		idx := len(s.BeforeDelete) + 1
		name := strs.Format(`%s_%d`, name, idx)
		Triggers[name] = trigger
		s.BeforeDelete = append(s.BeforeDelete, name)
	case AfterDelete:
		idx := len(s.AfterDelete) + 1
		name := strs.Format(`%s_%d`, name, idx)
		Triggers[name] = trigger
		s.AfterDelete = append(s.AfterDelete, name)
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

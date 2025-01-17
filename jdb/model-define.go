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

	idx := -1
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
	if col.Up() == KeyField.Up() {
		s.KeyField = col
	}
	if col.Up() == ProjectField.Up() {
		idx = slices.IndexFunc(s.Columns, func(e *Column) bool { return e == s.SourceField })
	}
	if idx == -1 {
		idx = slices.IndexFunc(s.Columns, func(e *Column) bool { return e == s.SystemKeyField })
	}
	if idx == -1 {
		idx = slices.IndexFunc(s.Columns, func(e *Column) bool { return e == s.IndexField })
	}
	if idx == -1 {
		s.Columns = append(s.Columns, col)
	} else {
		s.Columns = append(s.Columns[:idx], append([]*Column{col}, s.Columns[idx:]...)...)
	}
	if slices.Contains([]string{IndexField.Str(), ProjectField.Str(), CreatedAtField.Str(), UpdatedAtField.Str(), StateField.Str(), KeyField.Str(), SystemKeyField.Str(), SourceField.Str()}, name) {
		s.DefineIndex(true, name)
	} else if slices.Contains([]TypeData{TypeDataObject, TypeDataArray, TypeDataKey, TypeDataGeometry}, typeData) {
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
* DefineCreatedAtField
* @return *Column
**/
func (s *Model) DefineCreatedAtField() *Column {
	result := s.DefineColumn(CreatedAtField.Low(), CreatedAtField.TypeData())
	s.DefineIndex(true, CreatedAtField.Low())

	return result
}

/**
* DefineUpdatedAtField
* @return *Column
**/
func (s *Model) DefineUpdatedAtField() *Column {
	result := s.DefineColumn(UpdatedAtField.Low(), UpdatedAtField.TypeData())
	s.DefineIndex(true, UpdatedAtField.Low())

	return result
}

/**
* DefineStateField
* @return *Column
**/
func (s *Model) DefineStateField() *Column {
	result := s.DefineColumn(StateField.Low(), StateField.TypeData())
	s.DefineIndex(true, StateField.Low())

	return result
}

/**
* DefineKeyField
* @return *Column
**/
func (s *Model) DefineKeyField() *Column {
	result := s.DefineColumn(KeyField.Low(), KeyField.TypeData())
	s.DefineKey(KeyField.Low())

	return result
}

/**
* DefineSystemKeyField
* @return *Column
**/
func (s *Model) DefineSystemKeyField() *Column {
	result := s.DefineColumn(SystemKeyField.Low(), SystemKeyField.TypeData())
	s.DefineKey(SystemKeyField.Low())

	return result
}

/**
* DefineIndexField
* @return *Column
**/
func (s *Model) DefineIndexField() *Column {
	result := s.DefineColumn(IndexField.Low(), IndexField.TypeData())
	s.DefineIndex(true, IndexField.Low())

	return result
}

/**
* DefineClassField
* @return *Column
**/
func (s *Model) DefineClassField() *Column {
	result := s.DefineColumn(ClassField.Low(), ClassField.TypeData())
	result.Default = s.Low()
	s.DefineIndex(true, ClassField.Low())

	return result
}

/**
* DefineProjectField
* @return *Column
**/
func (s *Model) DefineProjectField() *Column {
	result := s.DefineColumn(ProjectField.Low(), ProjectField.TypeData())
	s.DefineIndex(true, ProjectField.Low())

	return result
}

/**
* DefineSourceField
* @return *Column
**/
func (s *Model) DefineSourceField() *Column {
	result := s.DefineColumn(SourceField.Low(), SourceField.TypeData())
	s.DefineIndex(true, SourceField.Low())

	return result
}

/**
* DefineGenerate
* @param name string
* @param function string
* @return *Model
**/
func (s *Model) DefineGenerated(name string, f FuncGenerated) {
	col := newColumn(s, name, "", TpGenerate, TypeDataNone, TypeDataNone.DefaultValue())
	col.FuncGenerated = f
	s.Columns = append(s.Columns, col)
}

/**
* DefineDetail
* @param name string
* @param detail Detail
* @return *Model
**/
func (s *Model) DefineDetail(name string) *Model {
	detail := NewModel(s.Schema, name, 1)
	col := newColumn(s, name, "", TpDetail, TypeDataNone, TypeDataNone.DefaultValue())
	col.Detail = detail
	s.Columns = append(s.Columns, col)
	s.Details[name] = detail
	detail.MakeDetailRelation(s)

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
func (s *Model) DefineDictionary(name, key, value string) *Dictionary {
	result := s.Dictionaries[value]
	if result != nil {
		return result
	}

	result = NewDictionary(s, name, key, value)
	s.Dictionaries[value] = result

	return result
}

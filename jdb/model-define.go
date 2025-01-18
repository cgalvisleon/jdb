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
	col := s.GetColumn(name)
	if col != nil {
		return col
	}

	idx := -1
	def := typeData.DefaultValue()
	col = newColumn(s, name, "", TpColumn, typeData, def)
	if col.Name == string(SourceField) {
		s.SourceField = col
	}
	if col.Name == string(CreatedAtField) {
		s.CreatedAtField = col
	}
	if col.Name == string(UpdatedAtField) {
		s.UpdatedAtField = col
	}
	if col.Name == string(SystemKeyField) {
		s.SystemKeyField = col
	}
	if col.Name == string(IndexField) {
		s.IndexField = col
	}
	if col.Name == string(StateField) {
		s.StateField = col
	}
	if col.Name == string(ClassField) {
		s.ClassField = col
	}
	if col.Name == string(KeyField) {
		s.KeyField = col
	}
	if col.Name == string(FullTextField) {
		s.FullTextField = col
	}
	if col.Name == string(ProjectField) {
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
	if slices.Contains([]string{string(IndexField), string(ProjectField), string(CreatedAtField), string(UpdatedAtField), string(StateField), string(KeyField), string(SystemKeyField), string(SourceField)}, name) {
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

	s.DefineColumn(string(SourceField), SourceField.TypeData())
	def := typeData.DefaultValue()
	col = newColumn(s, name, "", TpAtribute, typeData, def)
	col.Field = string(SourceField)
	s.Columns = append(s.Columns, col)

	return col
}

/**
* DefineCreatedAtField
* @return *Column
**/
func (s *Model) DefineCreatedAtField() *Column {
	result := s.DefineColumn(string(CreatedAtField), CreatedAtField.TypeData())
	s.DefineIndex(true, string(CreatedAtField))

	return result
}

/**
* DefineUpdatedAtField
* @return *Column
**/
func (s *Model) DefineUpdatedAtField() *Column {
	result := s.DefineColumn(string(UpdatedAtField), UpdatedAtField.TypeData())
	s.DefineIndex(true, string(UpdatedAtField))

	return result
}

/**
* DefineStateField
* @return *Column
**/
func (s *Model) DefineStateField() *Column {
	result := s.DefineColumn(string(StateField), StateField.TypeData())
	s.DefineIndex(true, string(StateField))

	return result
}

/**
* DefineKeyField
* @return *Column
**/
func (s *Model) DefineKeyField() *Column {
	result := s.DefineColumn(string(KeyField), KeyField.TypeData())
	s.DefineKey(string(KeyField))

	return result
}

/**
* DefineSystemKeyField
* @return *Column
**/
func (s *Model) DefineSystemKeyField() *Column {
	result := s.DefineColumn(string(SystemKeyField), SystemKeyField.TypeData())
	s.DefineKey(string(SystemKeyField))

	return result
}

/**
* DefineIndexField
* @return *Column
**/
func (s *Model) DefineIndexField() *Column {
	result := s.DefineColumn(string(IndexField), IndexField.TypeData())
	s.DefineIndex(true, string(IndexField))

	return result
}

/**
* DefineClassField
* @return *Column
**/
func (s *Model) DefineClassField() *Column {
	result := s.DefineColumn(string(ClassField), ClassField.TypeData())
	result.Default = s.Low()
	s.DefineIndex(true, string(ClassField))

	return result
}

/**
* DefineProjectField
* @return *Column
**/
func (s *Model) DefineProjectField() *Column {
	result := s.DefineColumn(string(ProjectField), ProjectField.TypeData())
	s.DefineIndex(true, string(ProjectField))

	return result
}

/**
* DefineSourceField
* @return *Column
**/
func (s *Model) DefineSourceField() *Column {
	result := s.DefineColumn(string(SourceField), SourceField.TypeData())
	s.DefineIndex(true, string(SourceField))

	return result
}

/**
* DefineFullText
* @param fields []string
* @return *Column
**/
func (s *Model) DefineFullText(fields []string) *Column {
	result := s.DefineColumn(string(FullTextField), FullTextField.TypeData())
	result.FullText = fields
	s.DefineIndex(true, string(FullTextField))

	return result
}

/**
* DefineGenerate
* @param name string
* @param function string
* @return *Model
**/
func (s *Model) DefineGenerated(name string, f FuncGenerated) {
	col := newColumn(s, name, "", TpGenerated, TypeDataNone, TypeDataNone.DefaultValue())
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
* DefineEvent
* @param tp TypeEvent
* @param event Event
* @return *Model
**/
func (s *Model) DefineEvent(tp TypeEvent, event Event) *Model {
	switch tp {
	case EventInsert:
		s.EventsInsert = append(s.EventsInsert, event)
	case EventUpdate:
		s.EventsUpdate = append(s.EventsUpdate, event)
	case EventDelete:
		s.EventsDelete = append(s.EventsDelete, event)
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

package jdb

import (
	"slices"
)

/**
* DefineIdxColumn
* @param name string, typeData TypeData
* @return *Column
**/
func (s *Model) DefineColumn(name string, typeData TypeData) *Column {
	col := s.GetColumn(name)
	if col != nil {
		return col
	}

	def := typeData.DefaultValue()
	col = newColumn(s, name, "", TpColumn, typeData, def)
	s.Columns = append(s.Columns, col)
	if slices.Contains([]string{string(IndexField), string(ProjectField), string(CreatedAtField), string(UpdatedAtField), string(StateField), string(PrimaryKeyField), string(SystemKeyField), string(SourceField)}, name) {
		s.DefineIndex(true, name)
	} else if slices.Contains([]TypeData{TypeDataObject, TypeDataArray, TypeDataKey, TypeDataGeometry, TypeDataTime}, typeData) {
		s.DefineIndex(true, name)
	}

	return col
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
		if col.TypeColumn == TpColumn {
			idx := NewIndex(col, sort)
			s.Indices = append(s.Indices, idx)
		}
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
		s.Uniques = append(s.Uniques, idx)
	}

	return s
}

/**
* DefineHidden
* @param colums ...string
* @return *Model
**/
func (s *Model) DefineHidden(colums ...string) *Model {
	cols := s.GetColumns(colums...)
	if len(cols) == 0 {
		return s
	}

	for _, col := range cols {
		col.Hidden = true
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
			s.ColRequired[col.Name] = true
		}
	}

	return s
}

/**
* DefinePrimaryKey
* @param name string
* @return *Column
**/
func (s *Model) DefinePrimaryKey(name string) *Column {
	result := s.DefineColumn(name, PrimaryKeyField.TypeData())
	idx := slices.IndexFunc(s.PrimaryKeys, func(e *Column) bool { return e == result })
	if idx != -1 {
		s.PrimaryKeys = append(s.PrimaryKeys, result)
	}

	return result
}

/**
* DefineKeyField
* @return *Column
**/
func (s *Model) DefinePrimaryKeyField() *Column {
	return s.DefinePrimaryKey(string(PrimaryKeyField))
}

/**
* DefineForeignKey
* @param name string, with *Model, pkn string
* @return *Column
**/
func (s *Model) DefineForeignKey(name string, with *Model) *Column {
	if len(with.PrimaryKeys) == 0 {
		return nil
	}

	pk := with.PrimaryKeys[0]
	result := s.DefineColumn(name, pk.TypeData)
	result.Detail = &Relation{
		With:            with,
		Fk:              pk,
		Limit:           -1,
		OnDeleteCascade: true,
		OnUpdateCascade: true,
	}
	s.DefineIndex(true, result.Name)
	idx := slices.IndexFunc(s.ForeignKeys, func(e *Column) bool { return e == result })
	if idx != -1 {
		s.ForeignKeys = append(s.ForeignKeys, result)
	}

	return result
}

/**
* DefineSourceField
* @return *Column
**/
func (s *Model) DefineSourceField(name string) *Column {
	result := s.DefineColumn(name, SourceField.TypeData())
	s.DefineIndex(true, name)

	return result
}

/**
* DefineAtribute
* @param name string
* @param typeData TypeData
* @param def interface{}
* @return *Model
**/
func (s *Model) DefineAtribute(name string, typeData TypeData) *Column {
	result := s.GetColumn(name)
	if result != nil {
		return result
	}

	if s.SourceField == nil {
		s.DefineSourceField(SOURCE)
	}

	def := typeData.DefaultValue()
	result = newColumn(s, name, "", TpAtribute, typeData, def)
	result.Source = s.SourceField
	s.Columns = append(s.Columns, result)

	return result
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
	s.StateField = result

	return result
}

/**
* DefineSystemKeyField
* @return *Column
**/
func (s *Model) DefineSystemKeyField() *Column {
	result := s.DefineColumn(string(SystemKeyField), SystemKeyField.TypeData())
	result.Hidden = true
	s.DefineIndex(true, string(SystemKeyField))
	s.SystemKeyField = result

	return result
}

/**
* DefineIndexField
* @return *Column
**/
func (s *Model) DefineIndexField() *Column {
	result := s.DefineColumn(string(IndexField), IndexField.TypeData())
	s.DefineIndex(true, string(IndexField))
	s.IndexField = result

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
* DefineFullText
* @param fields []string
* @return language string
* @return *Column
**/
func (s *Model) DefineFullText(language string, fields []string) *Column {
	cols := s.GetColumns(fields...)
	result := s.DefineColumn(string(FullTextField), FullTextField.TypeData())
	result.FullText = &FullText{
		Language: language,
		Columns:  cols,
	}
	result.Hidden = true
	s.DefineIndex(true, string(FullTextField))
	s.FullTextField = result

	return result
}

/**
* DefineGenerate
* @param name string, fn GeneratedFunction
* @return *Column
**/
func (s *Model) DefineGenerated(name string, fn GeneratedFunction) *Column {
	result := newColumn(s, name, "", TpGenerated, TypeDataNone, TypeDataNone.DefaultValue())
	result.GeneratedFunction = fn
	s.Columns = append(s.Columns, result)
	s.GeneratedFields = append(s.GeneratedFields, result)

	return result
}

/**
* DefineRelation
* @param name, relatedTo string
* @return *Relation
**/
func (s *Model) DefineRelation(name, relatedTo string) *Relation {
	if len(s.PrimaryKeys) == 0 {
		return nil
	}

	with := GetModel(relatedTo)
	if with == nil {
		with = NewModel(s.Schema, relatedTo, 0)
	}

	pk := s.PrimaryKeys[0]
	with.DefineAtribute(s.Name, pk.TypeData)
	with.DefineForeignKey(s.Name, s)
	col := newColumn(s, name, "", TpRelatedTo, TypeDataNone, TypeDataNone.DefaultValue())
	result := &Relation{
		With:  with,
		Fk:    pk,
		Limit: -1,
	}
	col.Detail = result
	s.Relations[name] = result

	return result
}

/**
* DefineDetail
* @param name, relatedTo string
* @return *Relation
**/
func (s *Model) DefineDetail(name string) *Relation {
	if len(s.PrimaryKeys) == 0 {
		return nil
	}

	relatedTo := s.Name + "_" + name
	with := GetModel(relatedTo)
	if with == nil {
		with = NewModel(s.Schema, relatedTo, 0)
	}

	pk := s.PrimaryKeys[0]
	with.DefineAtribute(s.Name, pk.TypeData)
	with.DefineForeignKey(s.Name, s)
	col := newColumn(s, name, "", TpRelatedTo, TypeDataNone, TypeDataNone.DefaultValue())
	result := &Relation{
		With:  with,
		Fk:    pk,
		Limit: -1,
	}
	col.Detail = result
	s.Details[name] = result

	return result
}

/**
* DefineHistory
* @param limit int64
* @return *Relation
**/
func (s *Model) DefineHistory(limit int64) *Relation {
	result := s.DefineDetail("historical")
	result.Limit = limit

	return result
}

/**
* DefineEvent
* @param tp TypeEvent
* @param event Event
* @return *Model
**/
func (s *Model) DefineEvent(tp TypeEvent, event Event) {
	switch tp {
	case EventInsert:
		s.EventsInsert = append(s.EventsInsert, event)
	case EventUpdate:
		s.EventsUpdate = append(s.EventsUpdate, event)
	case EventDelete:
		s.EventsDelete = append(s.EventsDelete, event)
	}
}

/**
* DefineModel
* @return *Model
**/
func (s *Model) DefineModel() *Model {
	s.DefineCreatedAtField()
	s.DefineUpdatedAtField()
	s.DefineStateField()
	s.DefineSystemKeyField()
	s.DefinePrimaryKey(PRIMARYKEY)
	s.DefineSourceField(SOURCE)
	s.DefineIndexField()

	return s
}

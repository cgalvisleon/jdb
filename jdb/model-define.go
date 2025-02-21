package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/strs"
)

/**
* DefineColumnIdx
* @param name string, typeData TypeData
* @return *Column
**/
func (s *Model) DefineColumnIdx(name string, typeData TypeData, idx int) *Column {
	col := s.GetColumn(name)
	if col != nil {
		return col
	}

	def := typeData.DefaultValue()
	col = newColumn(s, name, "", TpColumn, typeData, def)
	if idx == -1 {
		s.Columns = append(s.Columns, col)
	} else {
		s.Columns = append(s.Columns[:idx], append([]*Column{col}, s.Columns[idx:]...)...)
	}

	return col
}

/**
* DefineColumnIdx
* @param name string, typeData TypeData
* @return *Column
**/
func (s *Model) DefineColumn(name string, typeData TypeData) *Column {
	return s.DefineColumnIdx(name, typeData, -1)
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
			name := strs.Format("%s_%s_idx", s.Name, col.Name)
			s.Indices[name] = idx
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
		name := strs.Format("%s_%s_idx", s.Name, col.Name)
		s.Uniques[name] = idx
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
			s.Required[col.Name] = true
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
	nm := strs.Format("%s_%s_pk", s.Name, name)
	s.PrimaryKeys[nm] = result

	return result
}

/**
* DefineKeyField
* @return *Column
**/
func (s *Model) DefinePrimaryKeyField() *Column {
	return s.DefinePrimaryKey(PRIMARYKEY)
}

/**
* DefineForeignKey
* @param name string, with *Model, pkn string
* @return *Column
**/
func (s *Model) DefineForeignKey(name string, with *Model) *Column {
	pk := with.Pk()
	if pk == nil {
		return nil
	}

	result := s.DefineColumn(name, pk.TypeData)
	result.Detail = &Relation{
		With:            with,
		Fk:              pk,
		Limit:           -1,
		OnDeleteCascade: true,
		OnUpdateCascade: true,
	}
	s.DefineIndex(true, result.Name)
	nm := strs.Format("%s_%s_fk", s.Name, name)
	s.ForeignKeys[nm] = result

	return result
}

/**
* DefineSourceField
* @return *Column
**/
func (s *Model) DefineSource(name string) *Column {
	result := s.DefineColumn(name, SourceField.TypeData())
	s.DefineIndex(true, name)

	return result
}

/**
* DefineSourceField
* @return *Column
**/
func (s *Model) DefineSourceField() *Column {
	return s.DefineSource(SOURCE)
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
		s.DefineSourceField()
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
	idx := -1
	pk := s.Pk()
	if pk != nil {
		idx = slices.IndexFunc(s.Columns, func(e *Column) bool { return e == pk })
	}
	result := s.DefineColumnIdx(string(ProjectField), ProjectField.TypeData(), idx)
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
func (s *Model) DefineRelation(name, relatedTo, fkn string) *Relation {
	pk := s.Pk()
	if pk == nil {
		return nil
	}

	with := GetModel(relatedTo)
	if with == nil {
		with = NewModel(s.Schema, relatedTo, 1)
	}

	with.DefineAtribute(fkn, pk.TypeData)
	with.DefineForeignKey(fkn, s)
	col := newColumn(s, name, "", TpRelatedTo, TypeDataNone, TypeDataNone.DefaultValue())
	result := &Relation{
		With:  with,
		Fk:    pk,
		Limit: 0,
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
func (s *Model) DefineDetail(name, fkn string) *Model {
	relatedTo := s.Name + "_" + name
	result := s.DefineRelation(name, relatedTo, fkn)

	return result.With
}

/**
* DefineHistory
* @param limit int64
* @return *Relation
**/
func (s *Model) DefineHistory(limit int64) *Model {
	pk := s.Pk()
	if pk == nil {
		return nil
	}

	name := "historical"
	relatedTo := s.Name + "_" + name
	result := s.DefineRelation(name, relatedTo, pk.Name)
	result.Limit = limit
	result.With.DefineColumn(CREATED_AT, CreatedAtField.TypeData())
	result.With.DefineSourceField()
	result.With.DefineSystemKeyField()
	result.With.DefineColumn(HISTORY_INDEX, IndexField.TypeData())
	result.With.DefineIndex(true, HISTORY_INDEX)

	return result.With
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
* DefineEventError
* @param event Resilience
**/
func (s *Model) DefineEventError(event EventError) {
	s.EventError = event
}

/**
* DefineModel
* @return *Model
**/
func (s *Model) DefineModel() *Model {
	s.DefineCreatedAtField()
	s.DefineUpdatedAtField()
	s.DefineProjectField()
	s.DefineStateField()
	s.DefinePrimaryKeyField()
	s.DefineSourceField()
	s.DefineSystemKeyField()
	s.DefineIndexField()

	return s
}

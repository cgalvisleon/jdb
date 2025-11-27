package jdb

import (
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

/**
* defineColumn
* @param name string, typeData string, params et.Json
* @return error
**/
func (s *Model) defineColumn(name, typeData string, hidden bool, defaultValue interface{}) error {
	if !utility.ValidStr(name, 0, []string{}) {
		return fmt.Errorf(MSG_NAME_REQUIRED)
	}

	idx := s.getColumnIndex(name)
	if idx != -1 {
		return nil
	}

	if !TypeData[typeData] {
		return fmt.Errorf(MSG_TYPE_REQUIRED)
	}

	def := &Column{
		Name:    name,
		Type:    typeData,
		Default: defaultValue,
		Hidden:  hidden,
	}

	if def.Hidden {
		s.Hidden = append(s.Hidden, name)
	}

	if s.RecordField != "" {
		idx = s.getColumnIndex(s.RecordField)
		if idx != -1 {
			s.Columns = append(s.Columns[:idx], append([]*Column{def}, s.Columns[idx:]...)...)
			return nil
		}
	}

	s.Columns = append(s.Columns, def)
	return nil
}

/**
* defineColumns
* @param params et.Json
* @return error
**/
func (s *Model) defineColumns(params []et.Json) error {
	for _, param := range params {
		name := param.String("name")
		typeData := param.String("type")
		hidden := param.Bool("hidden")
		defaultValue := param["default"]
		err := s.defineColumn(name, typeData, hidden, defaultValue)
		if err != nil {
			return err
		}
	}

	return nil
}

/**
* DefineColumn
* @param name string, columnType string
* @return error
**/
func (s *Model) DefineColumn(name string, columnType string) error {
	return s.defineColumn(name, columnType, false, "")
}

/**
* DefineDefaulValue
* @param name string, defaultValue interface{}
* @return error
**/
func (s *Model) DefineDefaulValue(name string, defaultValue interface{}) error {
	idx := s.getColumnIndex(name)
	if idx == -1 {
		return nil
	}

	s.Columns[idx].Default = defaultValue
	return nil
}

/**
* DefineSourceField
* @param name string
* @return error
**/
func (s *Model) DefineSourceField(name string) error {
	if !utility.ValidStr(name, 0, []string{}) {
		return nil
	}

	err := s.defineColumn(name, TypeJson, false, "")
	if err != nil {
		return err
	}

	s.SourceField = name
	s.DefineIndexes(name)
	return nil
}

/**
* DefineRecordField
* @param name string
* @return error
**/
func (s *Model) DefineRecordField(name string) error {
	if !utility.ValidStr(name, 0, []string{}) {
		return nil
	}

	err := s.defineColumn(name, TypeKey, false, "")
	if err != nil {
		return err
	}

	s.RecordField = name
	s.DefineIndexes(name)
	return nil
}

/**
* DefineStatusField
* @param name string
* @return error
**/
func (s *Model) DefineStatusField(name string) error {
	if !utility.ValidStr(name, 0, []string{}) {
		return nil
	}

	err := s.defineColumn(name, TypeText, false, "")
	if err != nil {
		return err
	}

	s.StatusField = name
	s.DefineIndexes(name)
	return nil
}

/**
* DefineStatusFieldDefault
* @return error
**/
func (s *Model) DefineStatusFieldDefault() error {
	return s.DefineStatusField(STATUS)
}

/**
* DefineRecordFieldDefault
* @return error
**/
func (s *Model) DefineRecordFieldDefault() error {
	return s.DefineRecordField(RECORDID)
}

/**
* DefineSourceFieldDefault
* @return error
**/
func (s *Model) DefineSourceFieldDefault() error {
	return s.DefineSourceField(SOURCE)
}

/**
* DefineAtrib
* @param name string, defaultValue interface{}
* @return error
**/
func (s *Model) DefineAtrib(name string, defaultValue interface{}) error {
	if s.SourceField == "" {
		s.DefineSourceField(SOURCE)
	}

	return s.defineColumn(name, TypeAtribute, false, defaultValue)
}

/**
* DefineRequired
* @param names ...string
* @return
**/
func (s *Model) DefineRequired(names ...string) {
	for _, name := range names {
		idx := s.getColumnIndex(name)
		if idx == -1 {
			continue
		}

		s.Required = append(s.Required, name)
	}
}

/**
* DefineUniqueIndex
* @param names ...string
* @return
**/
func (s *Model) DefineUniqueIndex(names ...string) {
	for _, name := range names {
		idx := slices.Index(s.UniqueIndexes, name)
		if idx != -1 {
			continue
		}

		idx = s.getColumnIndex(name)
		if idx == -1 {
			continue
		}

		s.UniqueIndexes = append(s.UniqueIndexes, name)
	}
}

/**
* DefinePrimaryKeys
* @param names ...string
* @return
**/
func (s *Model) DefinePrimaryKeys(names ...string) {
	for _, name := range names {
		idx := slices.Index(s.PrimaryKeys, name)
		if idx != -1 {
			continue
		}

		idx = s.getColumnIndex(name)
		if idx == -1 {
			continue
		}

		s.DefineRequired(name)
		s.PrimaryKeys = append(s.PrimaryKeys, name)
	}
}

/**
* DefineForeignKey
* @param to *Model, fks []et.Json, onDelete string, onUpdate string
* @return error
**/
func (s *Model) DefineForeignKey(to *Model, fks []et.Json, onDelete string, onUpdate string) error {
	if utility.ValidStr(onDelete, 0, []string{}) && onDelete != "cascade" {
		return fmt.Errorf("on_delete must be cascade")
	}

	if utility.ValidStr(onUpdate, 0, []string{}) && onUpdate != "cascade" {
		return fmt.Errorf("on_update must be cascade")
	}

	s.ForeignKeys = append(s.ForeignKeys, et.Json{
		"schema": to.Schema,
		"name":   to.Name,
		"references": et.Json{
			"columns":   fks,
			"on_delete": onDelete,
			"on_update": onUpdate,
		},
	})
	return nil
}

/**
* DefineIndexes
* @param names ...string
* @return error
**/
func (s *Model) DefineIndexes(names ...string) error {
	for _, name := range names {
		idx := s.getColumnIndex(name)
		if idx == -1 {
			continue
		}

		idx = slices.Index(s.Indexes, name)
		if idx != -1 {
			continue
		}

		s.Indexes = append(s.Indexes, name)
	}

	return nil
}

/**
* DefineHidden
* @param names ...string
* @return
**/
func (s *Model) DefineHidden(names ...string) {
	for _, name := range names {
		idx := s.getColumnIndex(name)
		if idx == -1 {
			continue
		}

		s.Columns[idx].Hidden = true
		s.Hidden = append(s.Hidden, name)
	}
}

/**
* DefineDetail
* @param name string, fks et.Json, version int
* @return (*Model, error)
**/
func (s *Model) DefineDetail(name string, fks et.Json, version int) (*Model, error) {
	detailName := fmt.Sprintf("%s_%s", s.Name, name)
	result, _ := GetModel(s.Database, detailName)
	if result != nil {
		return result, nil
	}

	result, err := s.db.Define(et.Json{
		"schema":  s.Schema,
		"name":    detailName,
		"version": version,
	})
	if err != nil {
		return nil, err
	}

	for fk := range fks {
		err := result.defineColumn(fk, TypeKey, false, "")
		if err != nil {
			return nil, err
		}
	}

	err = result.DefineForeignKey(s, []et.Json{fks}, "cascade", "cascade")
	if err != nil {
		return nil, err
	}

	err = s.defineColumn(name, TypeDetail, false, "")
	if err != nil {
		return nil, err
	}

	s.Details[name] = &Detail{
		From:    result,
		Fks:     fks,
		Selects: []string{},
	}

	return result, nil
}

/**
* DefineRollup
* @param name string, from string, fks et.Json, selects []string
* @return *Model
**/
func (s *Model) DefineRollup(name, from string, fks et.Json, selects []string) error {
	model, err := GetModel(s.Database, from)
	if err != nil {
		return err
	}

	err = s.defineColumn(name, TypeRollup, false, "")
	if err != nil {
		return err
	}

	s.Rollups[name] = &Detail{
		From:    model,
		Fks:     fks,
		Selects: selects,
	}

	return nil
}

/**
* DefineRelation
* @param name string, from string, fks et.Json, selects []string
* @return error
**/
func (s *Model) DefineRelation(name, from string, fks et.Json, selects []string) error {
	model, err := GetModel(s.Database, from)
	if err != nil {
		return err
	}

	err = s.defineColumn(name, TypeRelation, false, "")
	if err != nil {
		return err
	}

	for fk := range fks {
		model.DefineIndexes(fk)
	}

	s.Relations[name] = &Detail{
		From:    model,
		Fks:     fks,
		Selects: selects,
	}

	return nil
}

/**
* DefineCalc
* @param name string, fn DataContext
* @return error
**/
func (s *Model) DefineCalc(name string, fn DataContext) error {
	err := s.defineColumn(name, TypeCalc, false, "")
	if err != nil {
		return err
	}

	s.Calcs[name] = fn
	return nil
}

/**
* DefineCreatedAtField
* @return *Model
**/
func (s *Model) DefineCreatedAtField() *Model {
	s.DefineColumn("created_at", TypeDateTime)
	return s
}

/**
* DefineUpdatedAtField
* @return *Model
**/
func (s *Model) DefineUpdatedAtField() *Model {
	s.DefineColumn("updated_at", TypeDateTime)
	return s
}

/**
* DefinePrimaryKeyField
* @return *Model
**/
func (s *Model) DefinePrimaryKeyField() *Model {
	s.DefineColumn(KEY, TypeKey)
	s.DefinePrimaryKeys(KEY)
	return s
}

/**
* DefineModel
* @return *Model
**/
func (s *Model) DefineModel() *Model {
	s.DefineCreatedAtField()
	s.DefineUpdatedAtField()
	s.DefineStatusFieldDefault()
	s.DefinePrimaryKeyField()
	s.DefineSourceFieldDefault()
	s.DefineRecordFieldDefault()
	return s
}

/**
* DefineTenantModel
* @return *Model
**/
func (s *Model) DefineProjectModel() *Model {
	s.DefineCreatedAtField()
	s.DefineUpdatedAtField()
	s.DefineStatusFieldDefault()
	s.DefinePrimaryKeyField()
	s.DefineColumn(PROJECT_ID, TypeKey)
	s.DefineSourceFieldDefault()
	s.DefineRecordFieldDefault()
	s.DefineIndexes(PROJECT_ID)
	return s
}

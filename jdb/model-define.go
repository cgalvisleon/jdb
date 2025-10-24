package jdb

import (
	"errors"
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

/**
* defineColumn
* @param name string, params et.Json
* @return error
**/
func (s *Model) defineColumn(name string, params et.Json) error {
	if !utility.ValidStr(name, 0, []string{}) {
		return fmt.Errorf(MSG_NAME_REQUIRED)
	}

	idx := s.getColumnIndex(name)
	if idx != -1 {
		return nil
	}

	typeData := params.String("type")
	if !TypeData[typeData] {
		return fmt.Errorf(MSG_TYPE_REQUIRED)
	}

	hidden := params.Bool("hidden")
	def := et.Json{
		"name":    name,
		"type":    typeData,
		"default": params.String("default"),
		"hidden":  hidden,
	}
	if s.RecordField != "" {
		idx = s.getColumnIndex(s.RecordField)
		if idx != -1 {
			s.Columns = append(s.Columns[:idx], append([]et.Json{def}, s.Columns[idx:]...)...)
			return nil
		}
	}

	s.Columns = append(s.Columns, def)
	if hidden {
		s.Hidden = append(s.Hidden, name)
	}
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
		err := s.defineColumn(name, param)
		if err != nil {
			return err
		}
	}

	return nil
}

/**
* defineRollup
* @param params et.Json
* @return error
**/
func (s *Model) defineRollup(params et.Json) error {
	name := params.String("name")
	if !utility.ValidStr(name, 0, []string{}) {
		return fmt.Errorf(MSG_ATRIB_REQUIRED, "name")
	}

	from := params.String("from")
	if !utility.ValidStr(from, 0, []string{}) {
		return fmt.Errorf(MSG_ATRIB_REQUIRED, "from")
	}

	model, err := s.GetModel(from)
	if err != nil {
		return err
	}

	as := "A"
	selectsOrigin := params.Json("selects")
	if selectsOrigin.IsEmpty() {
		return fmt.Errorf(MSG_ATRIB_REQUIRED, "selects")
	}

	selects := et.Json{}
	for k, v := range selectsOrigin {
		selects[fmt.Sprintf("%s.%s", as, k)] = v
	}

	fks := params.Json("fks")
	if fks.IsEmpty() {
		return fmt.Errorf(MSG_ATRIB_REQUIRED, "fks")
	}

	s.Rollups[name] = et.Json{
		"from": et.Json{
			model.Table: as,
		},
		"selects": selects,
		"fks":     fks,
	}
	return nil
}

/**
* DefineAtrib
* @param name string, defaultValue interface{}
* @return error
**/
func (s *Model) DefineAtrib(name string, defaultValue interface{}) error {
	if s.SourceField == "" {
		s.DefineSetSourceField(SOURCE)
	}

	return s.defineColumn(name, et.Json{
		"type":    TypeAtribute,
		"default": defaultValue,
	})
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

		s.Columns[idx].Set("hidden", true)
		s.Hidden = append(s.Hidden, name)
	}
}

/**
* DefineSetSourceField
* @param name string
* @return error
**/
func (s *Model) DefineSetSourceField(name string) error {
	if !utility.ValidStr(name, 0, []string{}) {
		return nil
	}

	SOURCE = name
	s.SourceField = name
	err := s.defineColumn(name, et.Json{
		"type": TypeJson,
	})
	if err != nil {
		return err
	}

	s.DefineIndexes(name)
	return nil
}

/**
* DefineSetRecordField
* @param name string
* @return error
**/
func (s *Model) DefineSetRecordField(name string) error {
	if !utility.ValidStr(name, 0, []string{}) {
		return nil
	}

	RECORDID = name
	s.RecordField = name
	err := s.defineColumn(name, et.Json{
		"type": TypeKey,
	})
	if err != nil {
		return err
	}

	s.DefineIndexes(name)
	return nil
}

/**
* DefineSetStatusField
* @param name string
* @return error
**/
func (s *Model) DefineSetStatusField(name string) error {
	if !utility.ValidStr(name, 0, []string{}) {
		return nil
	}

	STATUS = name
	s.StatusField = name
	err := s.defineColumn(name, et.Json{
		"type": TypeText,
	})
	if err != nil {
		return err
	}

	s.DefineIndexes(name)
	return nil
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
* DefineColumn
* @param name string, columnType string
* @return error
**/
func (s *Model) DefineColumn(name string, columnType string) error {
	return s.defineColumn(name, et.Json{
		"type": columnType,
	})
}

/**
* DefineRollups
* @param params et.Json
* @return error
**/
func (s *Model) DefineRollups(params []et.Json) error {
	for _, param := range params {
		err := s.defineRollup(param)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
* DefineRollup
* @param name string, from string, fks et.Json, selects et.Json
* @return error
**/
func (s *Model) DefineRollup(name, from string, fks, selects et.Json) error {
	err := s.defineRollup(et.Json{
		"name":    name,
		"from":    from,
		"selects": selects,
		"fks":     fks,
	})
	if err != nil {
		return err
	}

	return nil
}

/**
* DefineCalc
* @param name string, fn DataContext
* @return error
**/
func (s *Model) DefineCalc(name string, fn DataContext) error {
	err := s.defineColumn(name, et.Json{
		"type": TypeCalc,
	})
	if err != nil {
		return err
	}

	s.Calcs[name] = fn
	return nil
}

/**
* DefineColumnVm
* @param name string, script string
* @return error
**/
func (s *Model) DefineColumnVm(name string, script string) error {
	err := s.defineColumn(name, et.Json{
		"type": TypeVm,
	})
	if err != nil {
		return err
	}

	s.Vms[name] = script
	return nil
}

/**
* DefineDetail
* @param name string, fks []et.Json, version int
* @return (*Model, error)
**/
func (s *Model) DefineDetail(name string, fks []et.Json, version int) (*Model, error) {
	colName := name
	name = fmt.Sprintf("%s_%s", s.Name, name)
	result, err := s.GetModel(name)
	if err == nil {
		return result, nil
	}

	if !errors.Is(err, ErrModelNotFound) {
		return nil, err
	}

	result, err = s.db.Define(et.Json{
		"schema":  s.Schema,
		"name":    name,
		"version": version,
	})
	if err != nil {
		return nil, err
	}

	for _, fk := range fks {
		for f := range fk {
			result.defineColumn(f, et.Json{
				"type": TypeKey,
			})
		}
	}

	err = result.DefineForeignKey(s, fks, "cascade", "cascade")
	if err != nil {
		return nil, err
	}

	s.defineColumn(colName, et.Json{
		"type": TypeDetail,
	})
	s.details[colName] = result
	s.Details[colName] = result.ToJson()
	s.save()
	return result, nil
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
* DefineSourceField
* @return *Model
**/
func (s *Model) DefineSourceField() *Model {
	s.DefineSetSourceField(SOURCE)
	return s
}

/**
* DefineRecordField
* @return *Model
**/
func (s *Model) DefineRecordField() *Model {
	s.DefineSetRecordField(RECORDID)
	return s
}

/**
* DefineStatusField
* @return *Model
**/
func (s *Model) DefineStatusField() *Model {
	s.DefineSetStatusField(STATUS)
	return s
}

/**
* DefineModel
* @return *Model
**/
func (s *Model) DefineModel() *Model {
	s.DefineCreatedAtField()
	s.DefineUpdatedAtField()
	s.DefineStatusField()
	s.DefinePrimaryKeyField()
	s.DefineSourceField()
	s.DefineRecordField()
	return s
}

/**
* DefineTenantModel
* @return *Model
**/
func (s *Model) DefineTenantModel() *Model {
	s.DefineCreatedAtField()
	s.DefineUpdatedAtField()
	s.DefineStatusField()
	s.DefinePrimaryKeyField()
	s.DefineColumn(TENANT_ID, TypeKey)
	s.DefineSourceField()
	s.DefineRecordField()
	s.DefineIndexes(TENANT_ID)
	return s
}

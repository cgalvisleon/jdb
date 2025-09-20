package jdb

import (
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

/**
* save
* @return error
**/
func (s *Model) save() error {
	return setModel(s.Id, s.ToJson(), s.isDebug)
}

/**
* defineColumn
* @param name string, params et.Json
* @return error
**/
func (s *Model) defineColumn(name string, params et.Json) error {
	if !utility.ValidStr(name, 0, []string{}) {
		return fmt.Errorf(`%s (%s)`, MSG_NAME_REQUIRED, "defineColumn")
	}

	typeData := params.String("type")
	if !TypeData[typeData] {
		return fmt.Errorf(MSG_TYPE_REQUIRED)
	}

	s.Columns[name] = et.Json{
		"type":    typeData,
		"default": params.String("default"),
	}
	return nil
}

/**
* defineAtrib
* @param name string, defaultValue interface{}
* @return error
**/
func (s *Model) defineAtrib(name string, defaultValue interface{}) error {
	if s.SourceField == "" {
		s.defineSourceField(SOURCE)
	}

	return s.defineColumn(name, et.Json{
		"type":    TypeAtribute,
		"default": defaultValue,
	})
}

/**
* defineColumns
* @param params et.Json
* @return error
**/
func (s *Model) defineColumns(params et.Json) error {
	for k := range params {
		param := params.Json(k)
		err := s.defineColumn(k, param)
		if err != nil {
			return err
		}
	}

	return nil
}

/**
* definePrimaryKeys
* @param names ...string
* @return
**/
func (s *Model) definePrimaryKeys(names ...string) {
	for _, name := range names {
		idx := slices.Index(s.PrimaryKeys, name)
		if idx != -1 {
			continue
		}

		_, ok := s.Columns[name]
		if !ok {
			continue
		}

		s.defineRequired(name)
		s.PrimaryKeys = append(s.PrimaryKeys, name)
	}
}

/**
* defineIndices
* @param names ...string
* @return error
**/
func (s *Model) defineIndices(names ...string) error {
	for _, name := range names {
		idx := slices.Index(s.Indices, name)
		if idx != -1 {
			continue
		}

		_, ok := s.Columns[name]
		if !ok {
			continue
		}

		s.Indices = append(s.Indices, name)
	}

	return nil
}

/**
* defineRequired
* @param names ...string
* @return
**/
func (s *Model) defineRequired(names ...string) {
	for _, name := range names {
		idx := slices.Index(s.Required, name)
		if idx != -1 {
			continue
		}

		_, ok := s.Columns[name]
		if !ok {
			continue
		}

		s.Required = append(s.Required, name)
	}
}

/**
* defineSourceField
* @param name string
* @return error
**/
func (s *Model) defineSourceField(name string) error {
	if !utility.ValidStr(name, 0, []string{}) {
		return nil
	}

	s.SourceField = name
	err := s.defineColumn(name, et.Json{
		"type": TypeJson,
	})
	if err != nil {
		return err
	}

	s.defineIndices(name)
	return nil
}

/**
* defineRecordField
* @param name string
* @return error
**/
func (s *Model) defineRecordField(name string) error {
	if !utility.ValidStr(name, 0, []string{}) {
		return nil
	}

	s.RecordField = name
	err := s.defineColumn(name, et.Json{
		"type": TypeJson,
	})
	if err != nil {
		return err
	}

	s.defineIndices(name)
	return nil
}

/**
* defineForeignKeys
* @param params et.Json
* @return error
**/
func (s *Model) defineForeignKeys(params et.Json) error {
	for key := range params {
		param := params.Json(key)

		schema := param.String("schema")
		if !utility.ValidStr(schema, 0, []string{}) {
			return fmt.Errorf("schema is required")
		}

		name := param.String("name")
		if !utility.ValidStr(name, 0, []string{}) {
			return fmt.Errorf("name is required")
		}

		references := param.Json("references")
		if references.IsEmpty() {
			return fmt.Errorf("references is required")
		}

		columns := references.Json("columns")
		if columns.IsEmpty() {
			return fmt.Errorf("columns is required")
		}

		onDelete := references.String("on_delete")
		if utility.ValidStr(onDelete, 0, []string{}) && onDelete != "cascade" {
			return fmt.Errorf("on_delete must be cascade")
		}

		onUpdate := references.String("on_update")
		if utility.ValidStr(onUpdate, 0, []string{}) && onUpdate != "cascade" {
			return fmt.Errorf("on_update must be cascade")
		}

		model, err := s.db.getOrCreateModel(schema, name)
		if err != nil {
			return err
		}

		for key, val := range columns {
			pk, ok := model.GetColumn(val.(string))
			if !ok {
				return fmt.Errorf("column %s not found in %s", val, model.Id)
			}

			err := s.defineColumn(key, et.Json{
				"type": pk.String("type"),
			})
			if err != nil {
				return err
			}

			if !utility.ValidStr(onDelete, 0, []string{}) {
				continue
			}

			if !utility.ValidStr(onUpdate, 0, []string{}) {
				continue
			}

			err = s.defineIndices(key)
			if err != nil {
				return err
			}
		}

		s.ForeignKeys[key] = et.Json{
			"schema": schema,
			"name":   name,
			"references": et.Json{
				"columns":   columns,
				"on_delete": onDelete,
				"on_update": onUpdate,
			},
		}
	}

	return nil
}

/**
* defineDetails
* @param params et.Json
* @return error
**/
func (s *Model) defineDetails(params et.Json) error {
	for key := range params {
		param := params.Json(key)

		schema := param.String("schema")
		if !utility.ValidStr(schema, 0, []string{}) {
			return fmt.Errorf("schema is required")
		}

		name := param.String("name")
		if !utility.ValidStr(name, 0, []string{}) {
			return fmt.Errorf("name is required")
		}

		references := param.Json("references")
		if references.IsEmpty() {
			return fmt.Errorf("references is required")
		}

		columns := references.Json("columns")
		if columns.IsEmpty() {
			return fmt.Errorf("columns is required in references")
		}

		onDelete := references.String("on_delete")
		onUpdate := references.String("on_update")

		detail, err := s.db.getOrCreateModel(schema, name)
		if err != nil {
			return err
		}

		err = detail.defineForeignKeys(et.Json{
			s.Name: et.Json{
				"schema": s.Schema,
				"name":   s.Name,
				"references": et.Json{
					"columns":   columns,
					"on_delete": onDelete,
					"on_update": onUpdate,
				},
			},
		})
		if err != nil {
			return err
		}

		detail.masters[s.Name] = s
		detail.Masters[s.Name] = et.Json{
			"schema": s.Schema,
			"name":   s.Name,
			"references": et.Json{
				"columns": columns,
			},
		}

		err = detail.defineColumn(s.Name, et.Json{
			"type": TypeMaster,
		})
		if err != nil {
			return err
		}

		columns = param.Json("columns")
		err = detail.defineColumns(columns)
		if err != nil {
			return err
		}

		err = s.defineColumn(name, et.Json{
			"type": TypeDetail,
		})
		if err != nil {
			return err
		}

		s.details[name] = detail
		s.Details[key] = param
	}

	return nil
}

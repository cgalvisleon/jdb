package jdb

import (
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

/**
* DefineColumn
* @param name string, params et.Json
* @return error
**/
func (s *Model) DefineColumn(name string, params et.Json) error {
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

	s.Columns = append(s.Columns, et.Json{
		"name":    name,
		"type":    typeData,
		"default": params.String("default"),
	})
	return nil
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

	return s.DefineColumn(name, et.Json{
		"type":    TypeAtribute,
		"default": defaultValue,
	})
}

/**
* defineColumns
* @param params et.Json
* @return error
**/
func (s *Model) defineColumns(params []et.Json) error {
	for _, param := range params {
		name := param.String("name")
		err := s.DefineColumn(name, param)
		if err != nil {
			return err
		}
	}

	return nil
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
		idx := slices.Index(s.Indexes, name)
		if idx != -1 {
			continue
		}

		idx = s.getColumnIndex(name)
		if idx == -1 {
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
		idx := slices.Index(s.Required, name)
		if idx != -1 {
			continue
		}

		idx = s.getColumnIndex(name)
		if idx == -1 {
			continue
		}

		s.Required = append(s.Required, name)
	}
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

	SOURCE = name
	s.SourceField = name
	err := s.DefineColumn(name, et.Json{
		"type": TypeJson,
	})
	if err != nil {
		return err
	}

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

	RECORDID = name
	s.RecordField = name
	err := s.DefineColumn(name, et.Json{
		"type": TypeKey,
	})
	if err != nil {
		return err
	}

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

	STATUS = name
	s.StatusField = name
	err := s.DefineColumn(name, et.Json{
		"type": TypeJson,
	})
	if err != nil {
		return err
	}

	s.DefineIndexes(name)
	return nil
}

/**
* DefineForeignKeys
* @param params []et.Json
* @return error
**/
func (s *Model) DefineForeignKeys(params []et.Json) error {
	for _, param := range params {
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

		s.ForeignKeys = append(s.ForeignKeys, et.Json{
			"schema": schema,
			"name":   name,
			"references": et.Json{
				"columns":   columns,
				"on_delete": onDelete,
				"on_update": onUpdate,
			},
		})
	}

	return nil
}

/**
* DefineDetails
* @param params []et.Json
* @return error
**/
func (s *Model) DefineDetails(params []et.Json) error {
	for _, param := range params {
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

		columns := references.ArrayJson("columns")
		if len(columns) == 0 {
			return fmt.Errorf("columns is required in references")
		}

		onDelete := references.String("on_delete")
		onUpdate := references.String("on_update")

		detail, err := s.db.getOrCreateModel(schema, name)
		if err != nil {
			return err
		}

		err = detail.DefineForeignKeys([]et.Json{
			{
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

		err = detail.DefineColumn(s.Name, et.Json{
			"type": TypeMaster,
		})
		if err != nil {
			return err
		}

		err = detail.defineColumns(columns)
		if err != nil {
			return err
		}

		err = s.DefineColumn(name, et.Json{
			"type": TypeDetail,
		})
		if err != nil {
			return err
		}

		s.details[name] = detail
		s.Details[name] = param
	}

	return nil
}

/**
* DefineColumnCalc
* @param name string, fn DataContext
* @return error
**/
func (s *Model) DefineColumnCalc(name string, fn DataContext) error {
	err := s.DefineColumn(name, et.Json{
		"type": TypeCalc,
	})
	if err != nil {
		return err
	}

	s.Calcs[name] = fn
	return nil
}

/**
* DefineColumnRollup
* @param name string, script string
* @return error
**/
func (s *Model) DefineColumnVm(name string, script string) error {
	err := s.DefineColumn(name, et.Json{
		"type": TypeVm,
	})
	if err != nil {
		return err
	}

	s.Vms[name] = script
	return nil
}

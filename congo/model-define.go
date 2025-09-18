package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

/**
* defineColumn
* @param params et.Json
* @return error
**/
func (s *Model) defineColumn(params et.Json) error {
	name := params.String("name")
	if !utility.ValidStr(name, 0, []string{}) {
		return fmt.Errorf("name is required")
	}

	typeData := params.String("type")
	if !TypeData[typeData] {
		return fmt.Errorf("type is required")
	}

	s.Columns[name] = params
	return nil
}

/**
* defineAtrib
* @param name string
* @return error
**/
func (s *Model) defineAtrib(name string) error {
	if s.SourceField == "" {
		s.defineSourceField(SOURCE)
	}

	return s.defineColumn(et.Json{
		"name": name,
		"type": TypeAtribute,
	})
}

/**
* definePrimaryKeys
* @param names ...string
* @return
**/
func (s *Model) definePrimaryKeys(names ...string) {
	pks := []string{}
	for _, name := range names {
		_, ok := s.Columns[name]
		if !ok {
			continue
		}

		s.Required = append(s.Required, name)
		pks = append(pks, name)
	}

	if len(pks) == 0 {
		return
	}

	pk := fmt.Sprintf("pk_%s", s.Name)
	s.PrimaryKeys[pk] = pks
}

/**
* defineIndices
* @param names ...string
* @return error
**/
func (s *Model) defineIndices(names ...string) error {
	for _, name := range names {
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
	return s.defineColumn(et.Json{
		"name": name,
		"type": TypeJson,
	})
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

		model, err := s.db.getModel(schema, name)
		if err != nil {
			return err
		}

		for key, val := range columns {
			pk, ok := model.GetColumn(val.(string))
			if !ok {
				return fmt.Errorf("column %s not found in %s", val, model.Id)
			}

			err := s.defineColumn(et.Json{
				"name": key,
				"type": pk.String("type"),
			})
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
* defineRelations
* @param params et.Json
* @return error
**/
func (s *Model) defineRelations(params et.Json) error {
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

		relation, err := s.db.getModel(schema, name)
		if err != nil {
			return err
		}

		err = relation.defineForeignKeys(et.Json{
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

		s.Relations[key] = param
	}

	return nil
}

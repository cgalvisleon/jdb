package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

const (
	SOURCE        = "source"
	KEY           = "id"
	TypeInt       = "int"
	TypeFloat     = "float"
	TypeKey       = "key"
	TypeText      = "text"
	TypeMemo      = "memo"
	TypeDateTime  = "datetime"
	TypeBoolean   = "boolean"
	TypeJson      = "json"
	TypeSerial    = "serial"
	TypeBytes     = "bytes"
	TypeGeometry  = "geometry"
	TypeAtribute  = "atribute"
	TypeCalc      = "calc"
	TypeRelatedTo = "related_to"
	TypeRollup    = "rollup"
)

var (
	TypeData = map[string]bool{
		TypeInt:       true,
		TypeFloat:     true,
		TypeKey:       true,
		TypeText:      true,
		TypeMemo:      true,
		TypeDateTime:  true,
		TypeBoolean:   true,
		TypeJson:      true,
		TypeSerial:    true,
		TypeBytes:     true,
		TypeGeometry:  true,
		TypeAtribute:  true,
		TypeCalc:      true,
		TypeRelatedTo: true,
		TypeRollup:    true,
	}
)

type Model struct {
	Id          string    `json:"id"`
	Database    string    `json:"database"`
	Schema      string    `json:"schema"`
	Name        string    `json:"name"`
	Table       string    `json:"table"`
	Columns     et.Json   `json:"columns"`
	SourceField string    `json:"source_field"`
	Relations   et.Json   `json:"relations"`
	PrimaryKeys et.Json   `json:"primary_keys"`
	ForeignKeys et.Json   `json:"foreign_keys"`
	Indices     []string  `json:"indices"`
	Required    []string  `json:"required"`
	DDL         string    `json:"ddl"`
	db          *Database `json:"-"`
}

/**
* Define
* @param definition et.Json
* @return (*Model, error)
**/
func Define(definition et.Json) (*Model, error) {
	database := definition.String("database")
	schema := definition.String("schema")
	name := definition.String("name")
	result := &Model{
		Id:          fmt.Sprintf("%s.%s", schema, name),
		Database:    database,
		Schema:      schema,
		Name:        name,
		Table:       definition.String("table"),
		Columns:     et.Json{},
		SourceField: "",
		Relations:   et.Json{},
		PrimaryKeys: et.Json{},
		ForeignKeys: et.Json{},
		Indices:     definition.ArrayStr("indices"),
		Required:    definition.ArrayStr("required"),
	}

	columns := definition.Json("columns")
	for k := range columns {
		param := columns.Json(k)
		err := result.DefineColumn(param)
		if err != nil {
			return nil, err
		}
	}

	primaryKeys := definition.ArrayStr("primary_keys")
	result.DefinePrimaryKeys(primaryKeys...)

	sourceField := definition.String("source_field")
	err := result.DefineSourceField(sourceField)
	if err != nil {
		return nil, err
	}

	foreignKeys := definition.ArrayJson("foreign_keys")
	for _, foreignKey := range foreignKeys {
		err := result.DefineForeignKeys(foreignKey)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

/**
* ToJson
* @return et.Json
**/
func (s *Model) ToJson() et.Json {
	return et.Json{
		"id":           s.Id,
		"database":     s.Database,
		"schema":       s.Schema,
		"name":         s.Name,
		"table":        s.Table,
		"columns":      s.Columns,
		"source_field": s.SourceField,
		"relations":    s.Relations,
		"primary_keys": s.PrimaryKeys,
		"foreign_keys": s.ForeignKeys,
		"indices":      s.Indices,
		"required":     s.Required,
		"ddl":          s.DDL,
	}
}

/**
* DefineColumn
* @param params et.Json
* @return error
**/
func (s *Model) DefineColumn(params et.Json) error {
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
* DefineAtrib
* @param name string
* @return error
**/
func (s *Model) DefineAtrib(name string) error {
	if s.SourceField == "" {
		s.DefineSourceField(SOURCE)
	}

	return s.DefineColumn(et.Json{
		"name": name,
		"type": TypeAtribute,
	})
}

/**
* DefinePrimaryKeys
* @param names ...string
* @return
**/
func (s *Model) DefinePrimaryKeys(names ...string) {
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
* DefineIndices
* @param names ...string
* @return error
**/
func (s *Model) DefineIndices(names ...string) error {
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
* DefineRequired
* @param names ...string
* @return
**/
func (s *Model) DefineRequired(names ...string) {
	for _, name := range names {
		_, ok := s.Columns[name]
		if !ok {
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

	s.SourceField = name
	return s.DefineColumn(et.Json{
		"name": name,
		"type": TypeJson,
	})
}

/**
* DefineForeignKeys
* @param params et.Json
* @return error
**/
func (s *Model) DefineForeignKeys(params et.Json) error {
	columns := params.ArrayStr("columns")
	if len(columns) == 0 {
		return fmt.Errorf("columns is required")
	}

	references := params.Json("references")
	if references.IsEmpty() {
		return fmt.Errorf("references is required")
	}

	table := references.String("table")
	if !utility.ValidStr(table, 0, []string{}) {
		return fmt.Errorf("table is required in references")
	}

	referenceColumns := references.ArrayStr("columns")
	if len(referenceColumns) == 0 {
		return fmt.Errorf("columns is required in references")
	}

	onDelete := params.String("on_delete")
	if utility.ValidStr(onDelete, 0, []string{}) && onDelete != "cascade" {
		return fmt.Errorf("on_delete must be cascade")
	}

	onUpdate := params.String("on_update")
	if utility.ValidStr(onUpdate, 0, []string{}) && onUpdate != "cascade" {
		return fmt.Errorf("on_update must be cascade")
	}

	fk := fmt.Sprintf("fk_%s_%s", s.Name, table)
	s.ForeignKeys[fk] = et.Json{
		"columns": columns,
		"references": et.Json{
			"table":   table,
			"columns": referenceColumns,
		},
		"on_delete": onDelete,
		"on_update": onUpdate,
	}

	return nil
}

/**
* Load
* @return (string, error)
**/
func (s *Model) Load() (string, error) {
	if !utility.ValidStr(s.Database, 0, []string{}) {
		return "", fmt.Errorf("database is required")
	}

	if !utility.ValidStr(s.Schema, 0, []string{}) {
		return "", fmt.Errorf("schema is required")
	}

	if !utility.ValidStr(s.Name, 0, []string{}) {
		return "", fmt.Errorf("name is required")
	}

	if len(s.Columns) == 0 {
		err := s.DefineSourceField(SOURCE)
		if err != nil {
			return "", err
		}

		if len(s.PrimaryKeys) == 0 {
			s.DefineColumn(et.Json{
				"name": KEY,
				"type": TypeKey,
			})

			s.DefinePrimaryKeys(KEY)
		}
	}

	s.db = getDatabase(s.Database)
	return s.db.loadModel(s)
}

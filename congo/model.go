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

type DataFunctionTx func(tx *Tx, data et.Json) error

type Model struct {
	Id           string           `json:"id"`
	Database     string           `json:"database"`
	Schema       string           `json:"schema"`
	Name         string           `json:"name"`
	Table        string           `json:"table"`
	Columns      et.Json          `json:"columns"`
	SourceField  string           `json:"source_field"`
	Relations    et.Json          `json:"relations"`
	PrimaryKeys  et.Json          `json:"primary_keys"`
	ForeignKeys  et.Json          `json:"foreign_keys"`
	Indices      []string         `json:"indices"`
	Required     []string         `json:"required"`
	IsLocked     bool             `json:"is_locked"`
	db           *Database        `json:"-"`
	beforeInsert []DataFunctionTx `json:"-"`
	beforeUpdate []DataFunctionTx `json:"-"`
	beforeDelete []DataFunctionTx `json:"-"`
	afterInsert  []DataFunctionTx `json:"-"`
	afterUpdate  []DataFunctionTx `json:"-"`
	afterDelete  []DataFunctionTx `json:"-"`
}

/**
* Define
* @param definition et.Json
* @return (*Model, error)
**/
func Define(definition et.Json) (*Model, error) {
	database := definition.String("database")
	if !utility.ValidStr(database, 0, []string{}) {
		return nil, fmt.Errorf("database is required")
	}

	schema := definition.String("schema")
	if !utility.ValidStr(schema, 0, []string{}) {
		return nil, fmt.Errorf("schema is required")
	}

	name := definition.String("name")
	if !utility.ValidStr(name, 0, []string{}) {
		return nil, fmt.Errorf("name is required")
	}

	db := getDatabase(database)
	result, err := db.getModel(schema, name)
	result.Table = definition.String("table")
	result.Indices = definition.ArrayStr("indices")
	result.Required = definition.ArrayStr("required")

	columns := definition.Json("columns")
	for k := range columns {
		param := columns.Json(k)
		err := result.defineColumn(param)
		if err != nil {
			return nil, err
		}
	}

	primaryKeys := definition.ArrayStr("primary_keys")
	result.definePrimaryKeys(primaryKeys...)

	sourceField := definition.String("source_field")
	err = result.defineSourceField(sourceField)
	if err != nil {
		return nil, err
	}

	if err := result.validate(); err != nil {
		return nil, err
	}

	relations := definition.Json("relations")
	if !relations.IsEmpty() {
		err := result.defineRelations(relations)
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
		"is_locked":    s.IsLocked,
	}
}

/**
* validate
* @return error
**/
func (s *Model) validate() error {
	if len(s.Columns) != 0 {
		return nil
	}

	err := s.defineSourceField(SOURCE)
	if err != nil {
		return err
	}

	if len(s.PrimaryKeys) == 0 {
		s.defineColumn(et.Json{
			"name": KEY,
			"type": TypeKey,
		})

		s.definePrimaryKeys(KEY)
	}

	return nil
}

/**
* GetColumn
* @param name string
* @return (et.Json, bool)
**/
func (s *Model) GetColumn(name string) (et.Json, bool) {
	_, ok := s.Columns[name]
	if !ok {
		return et.Json{}, false
	}

	result := s.Columns.Json(name)
	return result, ok
}

/**
* Lock
* @return
**/
func (s *Model) Lock() {
	s.IsLocked = true
}

/**
* Unlock
* @return
**/
func (s *Model) Unlock() {
	s.IsLocked = false
}

/**
* Init
* @return error
**/
func (s *Model) Init() error {
	if !utility.ValidStr(s.Database, 0, []string{}) {
		return fmt.Errorf("database is required")
	}

	if !utility.ValidStr(s.Schema, 0, []string{}) {
		return fmt.Errorf("schema is required")
	}

	if !utility.ValidStr(s.Name, 0, []string{}) {
		return fmt.Errorf("name is required")
	}

	if err := s.validate(); err != nil {
		return err
	}

	return s.db.init(s)
}

package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
)

var dbs map[string]*Database

func init() {
	dbs = make(map[string]*Database)
}

type Database struct {
	Name   string            `json:"name"`
	Models map[string]*Model `json:"models"`
	Driver Driver            `json:"driver"`
}

/**
* ToJson
* @return et.Json
**/
func (s *Database) ToJson() et.Json {
	return et.Json{
		"name":   s.Name,
		"models": s.Models,
	}
}

/**
* GetDatabase
* @param name string
* @return *Database
**/
func getDatabase(name string) *Database {
	result, ok := dbs[name]
	if !ok {
		result = &Database{
			Name:   name,
			Models: make(map[string]*Model),
			Driver: drivers["postgres"],
		}
		dbs[name] = result
	}

	return result
}

/**
* Init
* @param model *Model
* @return error
**/
func (s *Database) init(model *Model) error {
	return s.Driver.Load(model)
}

/**
* getModel
* @param id string
* @return (*Model, error)
**/
func (s *Database) getModel(schema, name string) (*Model, error) {
	id := fmt.Sprintf("%s.%s", schema, name)
	model, ok := s.Models[id]
	if !ok {
		model = &Model{
			Id:           id,
			Database:     s.Name,
			Schema:       schema,
			Name:         name,
			Table:        "",
			Columns:      et.Json{},
			SourceField:  "",
			Relations:    et.Json{},
			PrimaryKeys:  et.Json{},
			ForeignKeys:  et.Json{},
			Indices:      []string{},
			Required:     []string{},
			beforeInsert: []DataFunctionTx{},
			beforeUpdate: []DataFunctionTx{},
			beforeDelete: []DataFunctionTx{},
			afterInsert:  []DataFunctionTx{},
			afterUpdate:  []DataFunctionTx{},
			afterDelete:  []DataFunctionTx{},
			db:           s,
		}
		s.Models[id] = model
	}

	return model, nil
}

/**
* GetModel
* @param database, schema, name string
* @return (*Model, error)
**/
func GetModel(database, schema, name string) (*Model, error) {
	db, ok := dbs[database]
	if !ok {
		return nil, fmt.Errorf("database %s not found", database)
	}

	id := fmt.Sprintf("%s.%s", schema, name)
	result, ok := db.Models[id]
	if !ok {
		return nil, fmt.Errorf("model %s not found", id)
	}

	return result, nil
}

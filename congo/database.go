package jdb

import (
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

var dbs map[string]*Database

func init() {
	dbs = make(map[string]*Database)
}

type Database struct {
	Name   string            `json:"name"`
	Models map[string]*Model `json:"models"`
	driver Driver            `json:"-"`
}

/**
* ToJson
* @return et.Json
**/
func (s *Database) ToJson() et.Json {
	bt, err := json.Marshal(s)
	if err != nil {
		return et.Json{}
	}

	var result et.Json
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return et.Json{}
	}

	return result
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
			driver: drivers[DriverPostgres],
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
	err := model.save()
	if err != nil {
		return err
	}

	return s.driver.Load(model)
}

/**
* getModel
* @param id string
* @return (*Model, error)
**/
func (s *Database) getModel(schema, name string) (*Model, error) {
	if !utility.ValidStr(schema, 0, []string{}) {
		return nil, fmt.Errorf(MSG_SCHEMA_REQUIRED)
	}

	if !utility.ValidStr(name, 0, []string{}) {
		return nil, fmt.Errorf(MSG_NAME_REQUIRED)
	}

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
			Details:      et.Json{},
			Masters:      et.Json{},
			Rollups:      et.Json{},
			PrimaryKeys:  et.Json{},
			ForeignKeys:  et.Json{},
			Indices:      []string{},
			Required:     []string{},
			db:           s,
			details:      make(map[string]*Model),
			masters:      make(map[string]*Model),
			rollups:      make(map[string]*Model),
			beforeInsert: []DataFunctionTx{},
			beforeUpdate: []DataFunctionTx{},
			beforeDelete: []DataFunctionTx{},
			afterInsert:  []DataFunctionTx{},
			afterUpdate:  []DataFunctionTx{},
			afterDelete:  []DataFunctionTx{},
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

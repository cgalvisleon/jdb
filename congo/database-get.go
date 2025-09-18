package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

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
			Atribs:       et.Json{},
			SourceField:  "",
			Details:      et.Json{},
			Masters:      et.Json{},
			Rollups:      et.Json{},
			PrimaryKeys:  []string{},
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
		model.BeforeInsert(model.beforeInsertDefault)
		model.BeforeUpdate(model.beforeUpdateDefault)
		model.BeforeDelete(model.beforeDeleteDefault)
		model.AfterInsert(model.afterInsertDefault)
		model.AfterUpdate(model.afterUpdateDefault)
		model.AfterDelete(model.afterDeleteDefault)
		s.Models[id] = model
	}

	return model, nil
}

/**
* GetDatabase
* @param name string
* @return (*Database, error)
**/
func GetDatabase(name string) (*Database, error) {
	result, ok := dbs[name]
	if !ok {
		return nil, fmt.Errorf(MSG_DATABASE_NOT_FOUND, name)
	}

	return result, nil
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

package jdb

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

var dbs map[string]*Database

func init() {
	dbs = make(map[string]*Database)
}

type Database struct {
	Name       string            `json:"name"`
	Models     map[string]*Model `json:"models"`
	UseCore    bool              `json:"use_core"`
	Connection et.Json           `json:"-"`
	driver     Driver            `json:"-"`
	db         *sql.DB           `json:"-"`
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
* @param name, driver string, params et.Json
* @return (*Database, error)
**/
func getDatabase(name, driver string, params et.Json) (*Database, error) {
	result, ok := dbs[name]
	if !ok {
		if _, ok := drivers[driver]; !ok {
			return nil, fmt.Errorf(MSG_DRIVER_NOT_FOUND, driver)
		}

		result = &Database{
			Name:       name,
			Models:     make(map[string]*Model),
			Connection: params,
		}
		result.driver = drivers[driver](result)
		err := result.load()
		if err != nil {
			return nil, err
		}

		dbs[name] = result
	}

	return result, nil
}

/**
* load
* @return error
**/
func (s *Database) load() error {
	if s.driver == nil {
		return fmt.Errorf(MSG_DRIVER_REQUIRED)
	}

	db, err := s.driver.Connect(s)
	if err != nil {
		return err
	}

	s.db = db

	if s.UseCore {
		err := initCore(s)
		if err != nil {
			console.Panic(err)
		}
	}

	return nil
}

/**
* initModel
* @param model *Model
* @return error
**/
func (s *Database) init(model *Model) error {
	err := s.driver.Load(model)
	if err != nil {
		return err
	}

	err = model.save()
	if err != nil {
		return err
	}

	return nil
}

/**
* getOrCreateModel
* @param schema, name string
* @return (*Model, error)
**/
func (s *Database) getOrCreateModel(schema, name string) (*Model, error) {
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
			Relations:    et.Json{},
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
* DefineModel
* @param definition et.Json
* @return (*Model, error)
**/
func (s *Database) DefineModel(definition et.Json) (*Model, error) {
	schema := definition.String("schema")
	if !utility.ValidStr(schema, 0, []string{}) {
		return nil, fmt.Errorf(MSG_SCHEMA_REQUIRED)
	}

	name := definition.String("name")
	if !utility.ValidStr(name, 0, []string{}) {
		return nil, fmt.Errorf(MSG_NAME_REQUIRED)
	}

	result, err := s.getOrCreateModel(schema, name)
	if err != nil {
		return nil, err
	}

	result.Table = definition.String("table")
	result.Indices = definition.ArrayStr("indices")
	result.Required = definition.ArrayStr("required")
	result.Version = definition.Int("version")

	columns := definition.Json("columns")
	err = result.defineColumns(columns)
	if err != nil {
		return nil, err
	}

	result.Atribs = definition.Json("atribs")
	for k, v := range result.Atribs {
		err := result.defineAtrib(k, v)
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

	details := definition.Json("details")
	if !details.IsEmpty() {
		err := result.defineDetails(details)
		if err != nil {
			return nil, err
		}
	}

	required := definition.ArrayStr("required")
	result.defineRequired(required...)

	debug := definition.Bool("debug")
	result.isDebug = debug

	return result, nil
}

/**
* getModel
* @param schema, name string
* @return (*Model, error)
**/
func (s *Database) getModel(schema, name string) (*Model, error) {
	id := fmt.Sprintf("%s.%s", schema, name)
	result, ok := s.Models[id]
	if !ok {
		return nil, fmt.Errorf(MSG_MODEL_NOT_FOUND, id)
	}

	return result, nil
}

/**
* query
* @param query *Ql
* @return (et.Items, error)
**/
func (s *Database) query(query *Ql) (et.Items, error) {
	if err := query.validate(); err != nil {
		return et.Items{}, err
	}

	result, err := s.driver.Query(query)
	if err != nil {
		return et.Items{}, err
	}

	if query.isDebug {
		console.Debugf("query:%s", query.ToJson().ToString())
	}

	return result, nil
}

/**
* exists
* @param query *Ql
* @return (bool, error)
**/
func (s *Database) exists(query *Ql) (bool, error) {
	if err := query.validate(); err != nil {
		return false, err
	}

	result, err := s.driver.Exists(query)
	if err != nil {
		return false, err
	}

	if query.isDebug {
		console.Debugf("exists:%s", query.ToJson().ToString())
	}

	return result, nil
}

/**
* count
* @param query *Ql
* @return (int, error)
**/
func (s *Database) count(query *Ql) (int, error) {
	if err := query.validate(); err != nil {
		return 0, err
	}

	result, err := s.driver.Count(query)
	if err != nil {
		return 0, err
	}

	if query.isDebug {
		console.Debugf("count:%s", query.ToJson().ToString())
	}

	return result, nil
}

/**
* command
* @param command *Command
* @return (et.Items, error)
**/
func (s *Database) command(command *Command) (et.Items, error) {
	if err := command.validate(); err != nil {
		return et.Items{}, err
	}

	result, err := s.driver.Command(command)
	if err != nil {
		return et.Items{}, err
	}

	if command.isDebug {
		console.Debugf("command:%s", command.toJson().ToString())
	}

	return result, nil
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

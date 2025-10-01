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
	Db         *sql.DB           `json:"-"`
	driver     Driver            `json:"-"`
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
* @param name, driver string, userCore bool, params et.Json
* @return (*Database, error)
**/
func getDatabase(name, driver string, userCore bool, params et.Json) (*Database, error) {
	result, ok := dbs[name]
	if !ok {
		if _, ok := drivers[driver]; !ok {
			return nil, fmt.Errorf(MSG_DRIVER_NOT_FOUND, driver)
		}

		result = &Database{
			Name:       name,
			Models:     make(map[string]*Model),
			UseCore:    userCore,
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

	s.Db = db
	if s.UseCore {
		err := initCore(s)
		if err != nil {
			console.Panic(err)
		}
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

	result, ok := s.Models[name]
	if !ok {
		result = &Model{
			Database:          s.Name,
			Schema:            schema,
			Name:              name,
			Table:             "",
			Columns:           []et.Json{},
			SourceField:       "",
			RecordField:       "",
			Details:           []et.Json{},
			Masters:           []et.Json{},
			Rollups:           []et.Json{},
			Relations:         []et.Json{},
			PrimaryKeys:       []string{},
			ForeignKeys:       []et.Json{},
			Indexes:           []string{},
			Required:          []string{},
			BeforeInserts:     make([]string, 0),
			BeforeUpdates:     make([]string, 0),
			BeforeDeletes:     make([]string, 0),
			AfterInserts:      make([]string, 0),
			AfterUpdates:      make([]string, 0),
			AfterDeletes:      make([]string, 0),
			db:                s,
			details:           make(map[string]*Model),
			masters:           make(map[string]*Model),
			eventBeforeInsert: []DataFunctionTx{},
			eventBeforeUpdate: []DataFunctionTx{},
			eventBeforeDelete: []DataFunctionTx{},
			eventAfterInsert:  []DataFunctionTx{},
			eventAfterUpdate:  []DataFunctionTx{},
			eventAfterDelete:  []DataFunctionTx{},
		}
		result.EventBeforeInsert(result.beforeInsertDefault)
		result.EventBeforeUpdate(result.beforeUpdateDefault)
		result.EventBeforeDelete(result.beforeDeleteDefault)
		result.EventAfterInsert(result.afterInsertDefault)
		result.EventAfterUpdate(result.afterUpdateDefault)
		result.EventAfterDelete(result.afterDeleteDefault)
		s.Models[name] = result
	}

	return result, nil
}

/**
* getModel
* @param name string
* @return (*Model, error)
**/
func (s *Database) getModel(name string) (*Model, error) {
	result, ok := s.Models[name]
	if ok {
		return result, nil
	}

	err := loadModel(name, &result)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf(MSG_MODEL_NOT_FOUND, name)
	}

	result.db = s
	result.EventBeforeInsert(result.beforeInsertDefault)
	result.EventBeforeUpdate(result.beforeUpdateDefault)
	result.EventBeforeDelete(result.beforeDeleteDefault)
	result.EventAfterInsert(result.afterInsertDefault)
	result.EventAfterUpdate(result.afterUpdateDefault)
	result.EventAfterDelete(result.afterDeleteDefault)
	s.Models[name] = result

	for _, item := range result.Details {
		name := item.String("name")
		detail, err := s.getModel(name)
		if err != nil {
			continue
		}

		if detail == nil {
			continue
		}

		result.details[name] = detail
		detail.masters[result.Name] = result
	}

	return result, nil
}

/**
* initModel
* @param model *Model
* @return error
**/
func (s *Database) init(model *Model) error {
	if err := model.prepare(); err != nil {
		return err
	}

	sql, err := s.driver.Load(model)
	if err != nil {
		return err
	}

	if model.isInit {
		return nil
	}

	if model.isDebug {
		console.Debugf("init:%s", model.ToJson().ToEscapeHTML())
	}

	_, err = s.Query(sql)
	if err != nil {
		return err
	}

	err = model.save()
	if err != nil {
		return err
	}

	model.SetInit()

	return nil
}

/**
* query
* @param query *Ql
* @return (et.Items, error)
**/
func (s *Database) query(query *Ql) (et.Items, error) {
	sql, err := s.driver.Query(query)
	if err != nil {
		return et.Items{}, err
	}

	if query.isDebug {
		console.Debugf("query:%s", query.ToJson().ToEscapeHTML())
	}

	result, err := s.QueryTx(query.tx, sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* command
* @param command *Cmd
* @return (et.Items, error)
**/
func (s *Database) command(command *Cmd) (et.Items, error) {
	sql, err := s.driver.Command(command)
	if err != nil {
		return et.Items{}, err
	}

	if command.isDebug {
		console.Debugf("command:%s", command.ToJson().ToEscapeHTML())
	}

	result, err := s.QueryTx(command.tx, sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* Define
* @param definition et.Json
* @return (*Model, error)
**/
func (s *Database) Define(definition et.Json) (*Model, error) {
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
	result.Indexes = definition.ArrayStr("indexes")
	result.Required = definition.ArrayStr("required")
	result.Version = definition.Int("version")

	columns := definition.ArrayJson("columns")
	err = result.defineColumns(columns)
	if err != nil {
		return nil, err
	}

	atribs := definition.Json("atribs")
	for k, v := range atribs {
		err := result.DefineAtrib(k, v)
		if err != nil {
			return nil, err
		}
	}

	primaryKeys := definition.ArrayStr("primary_keys")
	result.DefinePrimaryKeys(primaryKeys...)

	sourceField := definition.String("source_field")
	err = result.DefineSourceField(sourceField)
	if err != nil {
		return nil, err
	}

	statusField := definition.String("status_field")
	err = result.DefineStatusField(statusField)
	if err != nil {
		return nil, err
	}

	recordField := definition.String("record_field")
	err = result.DefineRecordField(recordField)
	if err != nil {
		return nil, err
	}

	details := definition.ArrayJson("details")
	if len(details) > 0 {
		err := result.DefineDetails(details)
		if err != nil {
			return nil, err
		}
	}

	required := definition.ArrayStr("required")
	result.DefineRequired(required...)

	debug := definition.Bool("debug")
	result.isDebug = debug

	return result, nil
}

/**
* Select
* @param query et.Json
* @return (*Ql, error)
**/
func (s *Database) Select(query et.Json) (*Ql, error) {
	result := newQl(s)
	from := query.Json("from")
	if from.IsEmpty() {
		return nil, fmt.Errorf(MSG_FROM_REQUIRED)
	}
	result.Froms = from

	return result.setQuery(query), nil
}

/**
* From
* @return *Ql
**/
func (s *Database) From(name string) *Ql {
	result := newQl(s)
	model, err := s.getModel(name)
	if err != nil {
		return result
	}

	result.addFrom(model, "A")
	return result
}

/**
* GetModel
* @param name string
* @return (*Model, error)
**/
func (s *Database) GetModel(name string) (*Model, error) {
	return s.getModel(name)
}

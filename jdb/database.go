package jdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

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

	id := fmt.Sprintf("%s.%s", schema, name)
	model, ok := s.Models[id]
	if !ok {
		model = &Model{
			Id:           id,
			Database:     s.Name,
			Schema:       schema,
			Name:         name,
			Table:        "",
			Columns:      []et.Json{},
			SourceField:  "",
			RecordField:  "",
			Details:      []et.Json{},
			Masters:      []et.Json{},
			Rollups:      []et.Json{},
			Relations:    []et.Json{},
			PrimaryKeys:  []string{},
			ForeignKeys:  []et.Json{},
			Indices:      []string{},
			Required:     []string{},
			db:           s,
			details:      make(map[string]*Model),
			masters:      make(map[string]*Model),
			calls:        make(map[string]*DataContext),
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
* initModel
* @param model *Model
* @return error
**/
func (s *Database) init(model *Model) error {
	if err := model.validate(); err != nil {
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
		console.Debug("load:\n\t", sql)
		return nil
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
	if err := query.validate(); err != nil {
		return et.Items{}, err
	}

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
	if err := command.validate(); err != nil {
		return et.Items{}, err
	}

	sql, err := s.driver.Command(command)
	if err != nil {
		return et.Items{}, err
	}

	if command.isDebug {
		console.Debugf("command:%s", command.toJson().ToString())
	}

	result, err := s.QueryTx(command.tx, sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* getModelByName
* @param name string
* @return (*Model, error)
**/
func (s *Database) getModelByName(name string) (*Model, error) {
	for _, model := range s.Models {
		if model.Name == name {
			return model, nil
		}
	}

	return nil, fmt.Errorf(MSG_MODEL_NOT_FOUND, name)
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
	result.Indices = definition.ArrayStr("indices")
	result.Required = definition.ArrayStr("required")
	result.Version = definition.Int("version")

	columns := definition.ArrayJson("columns")
	err = result.defineColumns(columns)
	if err != nil {
		return nil, err
	}

	atribs := definition.Json("atribs")
	for k, v := range atribs {
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

	recordField := definition.String("record_field")
	err = result.defineRecordField(recordField)
	if err != nil {
		return nil, err
	}

	if err := result.validate(); err != nil {
		return nil, err
	}

	details := definition.ArrayJson("details")
	if len(details) > 0 {
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
	lst := strings.Split(name, ".")
	if len(lst) == 2 {
		model, err := s.getModel(lst[0], lst[1])
		if err != nil {
			return result
		}
		result.addFrom(model.Table, "A")
	} else {
		for _, model := range s.Models {
			if model.Name == name {
				result.addFrom(model.Table, "A")
				break
			}
		}
	}

	return result
}

package jdb

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/utility"
)

var dbs map[string]*DB

func init() {
	dbs = make(map[string]*DB)
}

type DB struct {
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
func (s *DB) ToJson() et.Json {
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
* getDatabase
* @param name, driver string, userCore bool, params Connection
* @return (*DB, error)
**/
func getDatabase(name, driver string, userCore bool, params et.Json) (*DB, error) {
	result, ok := dbs[name]
	if !ok {
		if _, ok := drivers[driver]; !ok {
			return nil, fmt.Errorf(MSG_DRIVER_NOT_FOUND, driver)
		}

		result = &DB{
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
func (s *DB) load() error {
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
			logs.Panic(err)
		}
	}

	return nil
}

/**
* getOrCreateModel
* @param schema, name string
* @return (*Model, error)
**/
func (s *DB) getOrCreateModel(schema, name string) (*Model, error) {
	if !utility.ValidStr(schema, 0, []string{}) {
		return nil, fmt.Errorf(MSG_SCHEMA_REQUIRED)
	}

	if !utility.ValidStr(name, 0, []string{}) {
		return nil, fmt.Errorf(MSG_NAME_REQUIRED)
	}

	result, ok := s.Models[name]
	if !ok {
		result = &Model{
			Database:      s.Name,
			Schema:        schema,
			Name:          name,
			Table:         "",
			Columns:       []et.Json{},
			SourceField:   "",
			RecordField:   "",
			Details:       make(map[string]et.Json),
			Masters:       make(map[string]et.Json),
			Calcs:         make(map[string]DataContext),
			Vms:           make(map[string]string),
			Rollups:       make(map[string]et.Json),
			Relations:     make(map[string]et.Json),
			PrimaryKeys:   []string{},
			ForeignKeys:   []et.Json{},
			Indexes:       []string{},
			Required:      []string{},
			BeforeInserts: make([]string, 0),
			BeforeUpdates: make([]string, 0),
			BeforeDeletes: make([]string, 0),
			AfterInserts:  make([]string, 0),
			AfterUpdates:  make([]string, 0),
			AfterDeletes:  make([]string, 0),
			db:            s,
			details:       make(map[string]*Model),
			beforeInserts: []DataFunctionTx{},
			beforeUpdates: []DataFunctionTx{},
			beforeDeletes: []DataFunctionTx{},
			afterInserts:  []DataFunctionTx{},
			afterUpdates:  []DataFunctionTx{},
			afterDeletes:  []DataFunctionTx{},
		}
		result.BeforeInsert(result.beforeInsertDefault)
		result.BeforeUpdate(result.beforeUpdateDefault)
		result.BeforeDelete(result.beforeDeleteDefault)
		result.AfterInsert(result.afterInsertDefault)
		result.AfterUpdate(result.afterUpdateDefault)
		result.AfterDelete(result.afterDeleteDefault)
		s.Models[name] = result
	}

	return result, nil
}

/**
* getModel
* @param name string
* @return (*Model, error)
**/
func (s *DB) getModel(name string) (*Model, error) {
	result, ok := s.Models[name]
	if ok {
		return result, nil
	}

	err := loadModel(name, &result)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, ErrModelNotFound
	}

	result.Calcs = make(map[string]DataContext)
	result.db = s
	result.details = make(map[string]*Model)
	result.beforeInserts = []DataFunctionTx{}
	result.beforeUpdates = []DataFunctionTx{}
	result.beforeDeletes = []DataFunctionTx{}
	result.afterInserts = []DataFunctionTx{}
	result.afterUpdates = []DataFunctionTx{}
	result.afterDeletes = []DataFunctionTx{}
	result.BeforeInsert(result.beforeInsertDefault)
	result.BeforeUpdate(result.beforeUpdateDefault)
	result.BeforeDelete(result.beforeDeleteDefault)
	result.AfterInsert(result.afterInsertDefault)
	result.AfterUpdate(result.afterUpdateDefault)
	result.AfterDelete(result.afterDeleteDefault)

	for _, defDetail := range result.Details {
		detailName := defDetail.String("name")
		detail, err := s.getModel(detailName)
		if err != nil {
			continue
		}

		if detail == nil {
			continue
		}

		result.details[name] = detail
	}
	s.Models[name] = result

	return result, nil
}

/**
* initModel
* @param model *Model
* @return error
**/
func (s *DB) init(model *Model) error {
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

	if model.IsDebug {
		logs.Debugf("init:%s", model.ToJson().ToEscapeHTML())
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
func (s *DB) query(query *Ql) (et.Items, error) {
	sql, err := s.driver.Query(query)
	if err != nil {
		return et.Items{}, err
	}

	if query.IsDebug {
		logs.Debugf("query:%s", query.ToJson().ToEscapeHTML())
	}

	result, err := s.QueryTx(query.tx, sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* command
* @param cmd *Cmd
* @return (et.Items, error)
**/
func (s *DB) command(cmd *Cmd) (et.Items, error) {
	sql, err := s.driver.Command(cmd)
	if err != nil {
		return et.Items{}, err
	}

	if cmd.IsDebug {
		logs.Debugf("command:%s", cmd.ToJson().ToEscapeHTML())
	}

	result, err := s.QueryTx(cmd.tx, sql)
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
func (s *DB) Define(definition et.Json) (*Model, error) {
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
	err = result.DefineSetSourceField(sourceField)
	if err != nil {
		return nil, err
	}

	statusField := definition.String("status_field")
	err = result.DefineSetStatusField(statusField)
	if err != nil {
		return nil, err
	}

	recordField := definition.String("record_field")
	err = result.DefineSetRecordField(recordField)
	if err != nil {
		return nil, err
	}

	details := definition.ArrayJson("details")
	for _, detail := range details {
		detailName := detail.String("name")
		if !utility.ValidStr(detailName, 0, []string{}) {
			continue
		}

		fks := detail.ArrayJson("fks")
		if len(fks) == 0 {
			continue
		}

		version := detail.Int("version")
		_, err := result.DefineDetail(detailName, fks, version)
		if err != nil {
			return nil, err
		}
	}

	required := definition.ArrayStr("required")
	result.DefineRequired(required...)

	debug := definition.Bool("debug")
	result.IsDebug = debug

	return result, nil
}

/**
* Select
* @param query et.Json
* @return (*Ql, error)
**/
func (s *DB) Select(query et.Json) (*Ql, error) {
	result := newQl(nil, "")
	from := query.Json("from")
	if from.IsEmpty() {
		return nil, fmt.Errorf(MSG_FROM_REQUIRED)
	}
	result.Froms = from

	return result.setQuery(query), nil
}

/**
* From
* @param model *Model, as string
* @return *Ql
**/
func (s *DB) From(model *Model, as string) *Ql {
	result := newQl(model, as)
	result.IsDebug = model.IsDebug
	return result
}

/**
* GetModel
* @param name string
* @return (*Model, error)
**/
func (s *DB) GetModel(name string) (*Model, error) {
	return s.getModel(name)
}

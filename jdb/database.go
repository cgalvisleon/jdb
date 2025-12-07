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
			Columns:       make([]*Column, 0),
			SourceField:   "",
			RecordField:   "",
			Details:       make(map[string]*Detail),
			Rollups:       make(map[string]*Detail),
			Relations:     make(map[string]*Detail),
			Calcs:         make(map[string]DataContext),
			UniqueIndexes: []string{},
			PrimaryKeys:   []string{},
			ForeignKeys:   []et.Json{},
			Indexes:       []string{},
			Required:      []string{},
			db:            s,
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

	for k, item := range result.Details {
		detail, err := s.getModel(item.From.Name)
		if err != nil {
			continue
		}

		item.From = detail
		result.Details[k] = item
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

	model.isInit = true

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
	result.UniqueIndexes = definition.ArrayStr("unique_indexes")
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

	required := definition.ArrayStr("required")
	result.DefineRequired(required...)

	details := definition.Json("details")
	for detailName := range details {
		detail := details.Json(detailName)
		fks := detail.Json("fks")
		if len(fks) == 0 {
			continue
		}

		version := detail.Int("version")
		_, err := result.DefineDetail(detailName, fks, version)
		if err != nil {
			return nil, err
		}
	}

	rollups := definition.Json("rollups")
	for rollupName := range rollups {
		rollup := rollups.Json(rollupName)
		fks := rollup.Json("fks")
		if len(fks) == 0 {
			continue
		}

		from := rollup.String("from")
		selects := rollup.ArrayStr("selects")
		err := result.DefineRollup(rollupName, from, fks, selects)
		if err != nil {
			return nil, err
		}
	}

	relations := definition.Json("relations")
	for relationName := range relations {
		relation := relations.Json(relationName)
		fks := relation.Json("fks")
		if len(fks) == 0 {
			continue
		}

		from := relation.String("from")
		selects := relation.ArrayStr("selects")
		err := result.DefineRelation(relationName, from, fks, selects)
		if err != nil {
			return nil, err
		}
	}

	debug := definition.Bool("debug")
	result.IsDebug = debug

	return result, nil
}

/**
* GetModel
* @param name string
* @return (*Model, error)
**/
func (s *DB) GetModel(name string) (*Model, error) {
	return s.getModel(name)
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
* Select
* @param query et.Json
* @return *Ql
**/
func (s *DB) Select(query et.Json) (*Ql, error) {
	from := query.Json("from")
	if from.IsEmpty() {
		return nil, fmt.Errorf(MSG_FROM_REQUIRED)
	}
	var result *Ql
	for name := range from {
		as := from.String(name)
		model, err := s.getModel(name)
		if err != nil {
			return nil, err
		}

		if result == nil {
			result = newQl(model, as)
		} else {
			result.addFroms(model, as)
		}
	}

	result.setQuery(query)

	return result, nil
}

/**
* Data
* @param query et.Json
* @return *Ql
**/
func (s *DB) Data(query et.Json) (*Ql, error) {
	result, err := s.Select(query)
	if err != nil {
		return nil, err
	}
	result.IsDataSource = true
	return result, nil
}

package jdb

import (
	"database/sql"
	"encoding/json"
	"slices"
	"time"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/timezone"
)

type DB struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"update_at"`
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UseCore     bool      `json:"use_core"`
	NodeId      int       `json:"node_id"`
	driver      Driver    `json:"-"`
	schemas     []*Schema `json:"-"`
	models      []*Model  `json:"-"`
	isInit      bool      `json:"-"`
	db          *sql.DB   `json:"-"`
	IsDebug     bool      `json:"-"`
}

/**
* NewDatabaseById
* @param id, name, driver string
* @return *DB, error
**/
func NewDatabase(id, name, driver string) (*DB, error) {
	idx := slices.IndexFunc(conn.DBS, func(e *DB) bool { return e.Name == name })
	if idx != -1 {
		return conn.DBS[idx], nil
	}

	if driver == "" {
		return nil, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	if _, ok := conn.Drivers[driver]; !ok {
		return nil, mistake.Newf(MSG_DRIVER_NOT_FOUND, driver)
	}

	now := timezone.NowTime()
	result := &DB{
		CreatedAt: now,
		UpdateAt:  now,
		Id:        id,
		Name:      name,
		driver:    conn.Drivers[driver](),
		schemas:   make([]*Schema, 0),
		models:    make([]*Model, 0),
	}
	conn.DBS = append(conn.DBS, result)

	return result, nil
}

/**
* Serialize
* @return []byte, error
**/
func (s *DB) serialize() ([]byte, error) {
	result, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

/**
* Describe
* @return et.Json
**/
func (s *DB) Describe() et.Json {
	definition, err := s.serialize()
	if err != nil {
		return et.Json{}
	}

	result := et.Json{}
	err = json.Unmarshal(definition, &result)
	if err != nil {
		return et.Json{}
	}

	var schemas = make([]et.Json, 0)
	for _, schema := range s.schemas {
		schemas = append(schemas, schema.Describe())
	}

	result["kind"] = "db"
	result["schemas"] = schemas
	result["driver"] = s.driver.Name()

	return result
}

/**
* Load
* @param kind, name string, out interface{}
* @return error
**/
func (s *DB) Load(kind, name string, out interface{}) error {
	if !s.UseCore || !s.isInit {
		return nil
	}

	item, err := s.getModel(kind, name)
	if err != nil {
		return err
	}

	if !item.Ok {
		return mistake.Newf(MSG_MODEL_NOT_FOUND, name)
	}

	definition, err := item.Byte("definition")
	if err != nil {
		return err
	}

	if s.IsDebug {
		console.Debug(kind, ":", string(definition))
	}

	err = json.Unmarshal(definition, out)
	if err != nil {
		return err
	}

	return nil
}

/**
* Save
* @return error
**/
func (s *DB) Save() error {
	if !s.UseCore || !s.isInit {
		return nil
	}

	definition, err := s.serialize()
	if err != nil {
		return err
	}

	err = s.upsertModel("db", s.Name, 1, definition)
	if err != nil {
		return err
	}

	return nil
}

/**
* SetDebug
* @param debug bool
**/
func (s *DB) SetDebug(debug bool) {
	s.IsDebug = debug
}

/**
* Debug
**/
func (s *DB) Debug() {
	s.SetDebug(true)
}

/**
* Conected
* @param params *ConnectParams
* @return bool
**/
func (s *DB) Conected(params ConnectParams) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	if s.db != nil {
		return nil
	}

	db, err := s.driver.Connect(params)
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

/**
* Disconected
* @return error
**/
func (s *DB) Disconected() error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Disconnect()
}

/**
* AddSchema
* @param schema *Schema
**/
func (s *DB) addSchema(schema *Schema) {
	idx := slices.IndexFunc(s.schemas, func(e *Schema) bool { return e.Name == schema.Name })
	if idx != -1 {
		return
	}

	s.schemas = append(s.schemas, schema)
}

/**
* GetSchema
* @param name string
* @return *Schema
**/
func (s *DB) GetSchema(name string) *Schema {
	idx := slices.IndexFunc(s.schemas, func(e *Schema) bool { return e.Name == name })
	if idx != -1 {
		return s.schemas[idx]
	}

	return nil
}

/**
* GetModel
* @param name string
* @return *Model
**/
func (s *DB) GetModel(name string) *Model {
	idx := slices.IndexFunc(s.models, func(e *Model) bool { return e.Name == name })
	if idx != -1 {
		return s.models[idx]
	}

	return nil
}

/**
* DropSchema
* @param name string
* @return error
**/
func (s *DB) DropSchema(name string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	schema := s.GetSchema(name)
	if schema == nil {
		return mistake.Newf(MSG_SCHEMA_NOT_FOUND, name)
	}

	for _, model := range schema.models {
		err := s.DropModel(model)
		if err != nil {
			return err
		}
	}

	err := s.deleteModel("schema", name)
	if err != nil {
		return err
	}

	return nil
}

/**
* LoadModel
* @param model *Model
* @return error
**/
func (s *DB) LoadModel(model *Model) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.LoadModel(model)
}

/**
* MutateModel
* @param model *Model
* @return error
**/
func (s *DB) MutateModel(model *Model) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.MutateModel(model)
}

/**
* DropModel
* @param model *Model
**/
func (s *DB) DropModel(model *Model) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	err := s.driver.DropModel(model)
	if err != nil {
		return err
	}

	schema := model.schema
	if schema != nil {
		schema.dropModel(model)
	}

	err = s.deleteModel("model", model.Name)
	if err != nil {
		return err
	}

	return nil
}

/**
* EmptyModel
* @param model *Model
* @return error
**/
func (s *DB) EmptyModel(model *Model) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.EmptyModel(model)
}

/**
* From
* @param name string
* @return *Ql
**/
func (s *DB) From(name string) *Ql {
	model := s.GetModel(name)
	if model == nil {
		return nil
	}

	return From(model)
}

/**
* Select
* @param ql *Ql
* @return et.Items, error
**/
func (s *DB) Select(ql *Ql) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Select(ql)
}

/**
* Count
* @param ql *Ql
* @return int, error
**/
func (s *DB) Count(ql *Ql) (int, error) {
	if s.driver == nil {
		return 0, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Count(ql)
}

/**
* Exists
* @param ql *Ql
* @return bool, error
**/
func (s *DB) Exists(ql *Ql) (bool, error) {
	if s.driver == nil {
		return false, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Exists(ql)
}

/**
* Command
* @param command *Command
* @return et.Item, error
**/
func (s *DB) Command(command *Command) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Command(command)
}

/**
* Query
* @param sql string, arg ...any
* @return et.Items, error
**/
func (s *DB) Query(sql string, arg ...any) (et.Items, error) {
	if s.db == nil {
		return et.Items{}, mistake.New(MSG_DATABASE_NOT_CONNECTED)
	}

	return Query(s.db, sql, arg...)
}

/**
* One
* @param sql string, arg ...any
* @return et.Item, error
**/
func (s *DB) One(sql string, arg ...any) (et.Item, error) {
	result, err := s.Query(sql, arg...)
	if err != nil {
		return et.Item{}, err
	}

	return result.First(), nil
}

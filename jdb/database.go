package jdb

import (
	"encoding/json"
	"slices"
	"strings"
	"time"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/timezone"
)

func Name(name string) string {
	name = strings.ToLower(name)
	return strs.ReplaceAll(name, []string{" "}, "_")
}

type DB struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"update_at"`
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UseCore     bool      `json:"use_core"`
	driver      Driver    `json:"-"`
	schemas     []*Schema `json:"-"`
	models      []*Model  `json:"-"`
	isInit      bool      `json:"-"`
	IsDebug     bool      `json:"-"`
}

/**
* NewDatabase
* @param name, driver string
* @return *DB
**/
func NewDatabase(name, driver string) (*DB, error) {
	if driver == "" {
		return nil, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	if _, ok := conn.Drivers[driver]; !ok {
		return nil, mistake.Newf(MSG_DRIVER_NOT_FOUND, driver)
	}

	name = Name(name)
	idx := slices.IndexFunc(conn.DBS, func(e *DB) bool { return e.Name == name })
	if idx != -1 {
		return conn.DBS[idx], nil
	}

	now := timezone.NowTime()
	result := &DB{
		CreatedAt: now,
		UpdateAt:  now,
		Id:        reg.GenId("db"),
		Name:      name,
		UseCore:   false,
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
* Debug
**/
func (s *DB) Debug() {
	s.IsDebug = true
}

/**
* Conected
* @return bool
**/
func (s *DB) Conected(params et.Json) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Connect(params)
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
* getTableModel
* @param name string
* @return *Model, error
**/
func (s *DB) getTableModel(name string) (*Model, error) {
	list := strs.Split(name, ".")
	if len(list) != 2 {
		return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, name)
	}

	model := s.GetModel(list[1])
	if model != nil {
		return model, nil
	}

	schema := s.GetSchema(list[0])
	if schema == nil {
		schema = NewSchema(s, list[0])
	}

	model = NewModel(schema, list[1], 1)
	if err := model.Init(); err != nil {
		return nil, err
	}

	return model, nil
}

/**
* SetMain
* @param params et.Json
* @return error
**/
func (s *DB) SetMain(params et.Json) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.SetMain(params)
}

/**
* GrantPrivileges
* @return error
**/
func (s *DB) GrantPrivileges(username, database string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.GrantPrivileges(username, database)
}

/**
* CreateUser
* @return error
**/
func (s *DB) CreateUser(username, password, confirmation string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.CreateUser(username, password, confirmation)
}

/**
* ChangePassword
* @return error
**/
func (s *DB) ChangePassword(username, password, confirmation string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.ChangePassword(username, password, confirmation)
}

/**
* DeleteUser
* @return error
**/
func (s *DB) DeleteUser(username string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.DeleteUser(username)
}

/**
* LoadSchema
* @param name string
* @return error
**/
func (s *DB) LoadSchema(name string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.LoadSchema(name)
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

	err := s.driver.DropSchema(name)
	if err != nil {
		return err
	}

	err = s.deleteModel("schema", name)
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

	schema := s.GetSchema(model.Schema)
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
* Sync
* @param command string, data et.Json
* @return error
**/
func (s *DB) Sync(command string, data et.Json) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Sync(command, data)
}

/**
* QueryTx
* @param tx *Tx, sql string, arg ...any
* @return et.Items, error
**/
func (s *DB) QueryTx(tx *Tx, sql string, arg ...any) (et.Items, error) {
	return s.driver.QueryTx(tx, sql, arg...)
}

/**
* Query
* @param sql string, arg ...any
* @return et.Items, error
**/
func (s *DB) Query(sql string, arg ...any) (et.Items, error) {
	return s.driver.Query(sql, arg...)
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

/**
* From
* @param name string
* @return *Ql, error
**/
func (s *DB) From(name string) (*Ql, error) {
	model, err := s.getTableModel(name)
	if err != nil {
		return nil, err
	}

	result := From(model)

	return result, nil
}

/**
* EventSync
**/
func (s *DB) EventSync() {
	syncChannel := strs.Format("sync:%s", s.Name)
	event.Subscribe(syncChannel, func(msg event.Message) {
		data := msg.Data
		fromId := data.ValStr("", "fromId")
		if fromId == "" || fromId == s.Id {
			return
		}

		command := data.ValStr("", "command")
		if command == "" {
			return
		}

		switch command {
		case "insert":
			s.Sync("insert", data)
		case "update":
			s.Sync("update", data)
		case "delete":
			s.Sync("delete", data)
		}
	})
}

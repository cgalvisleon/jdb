package jdb

import (
	"encoding/json"
	"slices"
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/strs"
)

func Name(name string) string {
	return strs.ReplaceAll(name, []string{" "}, "_")
}

type Mode int

const (
	Origin Mode = iota
	Replica
)

type DB struct {
	CreatedAt   time.Time `json:"created_date"`
	UpdateAt    time.Time `json:"update_date"`
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UseCore     bool      `json:"use_core"`
	Mode        Mode      `json:"mode"`
	schemas     []*Schema `json:"-"`
	driver      Driver    `json:"-"`
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

	if _, ok := Jdb.Drivers[driver]; !ok {
		return nil, mistake.Newf(MSG_DRIVER_NOT_FOUND, driver)
	}

	name = Name(name)
	idx := slices.IndexFunc(Jdb.DBS, func(db *DB) bool { return db.Name == name })
	if idx != -1 {
		return Jdb.DBS[idx], nil
	}

	now := time.Now()
	result := &DB{
		CreatedAt: now,
		UpdateAt:  now,
		Id:        reg.Id("db"),
		Name:      name,
		UseCore:   false,
		Mode:      Origin,
		schemas:   make([]*Schema, 0),
		driver:    Jdb.Drivers[driver](),
	}
	Jdb.DBS = append(Jdb.DBS, result)

	return result, nil
}

/**
* Load
* @return error
**/
func (s *DB) Load(kind, name string, out interface{}) error {
	if !s.UseCore {
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
	if !s.UseCore {
		return nil
	}

	buf, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = s.upsertModel("db", s.Name, 1, buf)
	if err != nil {
		return err
	}

	return nil
}

/**
* Describe
* @return et.Json
**/
func (s *DB) Describe() et.Json {
	var schemas = make([]et.Json, 0)
	for _, schema := range s.schemas {
		schemas = append(schemas, schema.Describe())
	}

	return et.Json{
		"created_date": s.CreatedAt,
		"update_date":  s.UpdateAt,
		"id":           s.Id,
		"name":         s.Name,
		"description":  s.Description,
		"schemas":      schemas,
	}
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
* GetSchema
* @param name string
* @return *Schema
**/
func (s *DB) GetSchema(name string) *Schema {
	idx := slices.IndexFunc(s.schemas, func(schema *Schema) bool { return schema.Name == name })
	if idx != -1 {
		return s.schemas[idx]
	}

	return nil
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
* CreateSchema
* @param name string
* @return error
**/
func (s *DB) CreateSchema(name string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.CreateSchema(name)
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

	return s.driver.DropSchema(name)
}

/**
* LoadModel
* @param model *Model
* @return error
**/
func (s *DB) LoadModel(model *Model) (bool, error) {
	if s.driver == nil {
		return false, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.LoadModel(model)
}

/**
* CreateModel
* @param model *Model
* @return error
**/
func (s *DB) CreateModel(model *Model) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.CreateModel(model)
}

/**
* DropModel
* @param model *Model
**/
func (s *DB) DropModel(model *Model) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.DropModel(model)
}

/**
* SaveModel
* @param model *Model
* @return error
**/
func (s *DB) SaveModel(model *Model) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.SaveModel(model)
}

/**
* Exec
* @param sql string
* @param params ...any
* @return error
**/
func (s *DB) Exec(sql string, params ...any) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Exec(sql, params...)
}

/**
* Query
* @param sql string
* @param params ...any
* @return et.Items, error
**/
func (s *DB) Query(sql string, params ...any) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Query(sql, params...)
}

/**
* One
* @param sql string
* @param params ...any
* @return et.Item, error
**/
func (s *DB) One(sql string, params ...any) (et.Item, error) {
	if s.driver == nil {
		return et.Item{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.One(sql, params...)
}

/**
* Data
* @param source string
* @param sql string
* @param params ...any
* @return et.Item, error
**/
func (s *DB) Data(source, sql string, params ...any) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Data(source, sql, params...)
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
* EventSync
**/
func (s *DB) EventSync() {
	syncChannel := strs.Format("sync:%s", s.Name)
	event.Subscribe(syncChannel, func(msg event.EvenMessage) {
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

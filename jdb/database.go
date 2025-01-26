package jdb

import (
	"time"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

func Name(name string) string {
	result := strs.ReplaceAll(name, []string{" "}, "_")
	return strs.Lowcase(result)
}

var JDBS []*DB = []*DB{}

type DB struct {
	CreatedAt   time.Time          `json:"created_date"`
	UpdateAt    time.Time          `json:"update_date"`
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Schemas     map[string]*Schema `json:"schemas"`
	UseCore     bool               `json:"use_core"`
	Node        int64              `json:"node"`
	driver      Driver
}

/**
* NewDatabase
* @param driver string
* @return *DB
**/
func NewDatabase(name, driver string) (*DB, error) {
	if driver == "" {
		return nil, console.Alertm(MSG_DRIVER_NOT_DEFINED)
	}

	if _, ok := drivers[driver]; !ok {
		return nil, console.Alertf(MSG_DRIVER_NOT_FOUND, driver)
	}

	now := time.Now()
	result := &DB{
		CreatedAt:   now,
		UpdateAt:    now,
		Id:          utility.UUID(),
		Name:        Name(name),
		Description: "",
		Schemas:     map[string]*Schema{},
		Node:        envar.GetInt64(1, "DB_NODE"),
		driver:      drivers[driver](),
	}
	utility.SetSnowflakeNode(result.Node)
	JDBS = append(JDBS, result)

	return result, nil
}

/**
* Describe
* @return et.Json
**/
func (s *DB) Describe() et.Json {
	result, err := et.Object(s)
	if err != nil {
		return et.Json{}
	}

	return result
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
func (s *DB) GrantPrivileges(username, dbName string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.GrantPrivileges(username, dbName)
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
* CreateCore
* @return error
**/
func (s *DB) CreateCore() error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.CreateCore()
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

	return s.driver.DropModel(model)
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
* All
* @param sql string
* @param params ...any
* @return et.Items, error
**/
func (s *DB) All(tp TypeSelect, sql string, params ...any) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.All(tp, sql, params...)
}

/**
* One
* @param sql string
* @param params ...any
* @return et.Item, error
**/
func (s *DB) One(tp TypeSelect, sql string, params ...any) (et.Item, error) {
	if s.driver == nil {
		return et.Item{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.One(tp, sql, params...)
}

/**
* Query
* @param ql *Ql
* @return et.Items, error
**/
func (s *DB) Query(ql *Ql) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Query(ql)
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
* GetSerie
* @return int64
**/
func (s *DB) GetSerie(tag string) int64 {
	if s.driver == nil {
		return 0
	}

	return s.driver.GetSerie(tag)
}

/**
* NextCode
* @param tag string
* @param format string "%08v" "PREFIX-%08v-SUFFIX"
* @return string
**/
func (s *DB) NextCode(tag, format string) string {
	if s.driver == nil {
		return ""
	}

	return s.driver.NextCode(tag, format)
}

/**
* SetSerie
* @return int64
**/
func (s *DB) SetSerie(tag string, val int) int64 {
	if s.driver == nil {
		return 0
	}

	return s.driver.SetSerie(tag, val)
}

/**
* CurrentSerie
* @param tag string
* @return int64
**/
func (s *DB) CurrentSerie(tag string) int64 {
	if s.driver == nil {
		return 0
	}

	return s.driver.CurrentSerie(tag)
}

/**
* SetKey
* @param key string
* @param value value []byte
* @return error
**/
func (s *DB) SetKey(key string, value []byte) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.SetKey(key, value)
}

/**
* GetKey
* @param key string
* @return et.KeyValue, error
**/
func (s *DB) GetKey(key string) (et.KeyValue, error) {
	if s.driver == nil {
		return et.KeyValue{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.GetKey(key)
}

/**
* DeleteKey
* @param key string
* @return error
**/
func (s *DB) DeleteKey(key string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.DeleteKey(key)
}

/**
* FindKeys
* @param search string
* @param page int
* @param rows int
* @return et.List, error
**/
func (s *DB) FindKeys(search string, page, rows int) (et.List, error) {
	if s.driver == nil {
		return et.List{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.FindKeys(search, page, rows)
}

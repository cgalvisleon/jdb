package jdb

import (
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/utility"
)

type DB struct {
	CreatedAt   time.Time          `json:"created_date"`
	UpdateAt    time.Time          `json:"update_date"`
	Id          string             `json:"id"`
	Description string             `json:"description"`
	Schemas     map[string]*Schema `json:"schemas"`
	driver      Driver
}

/**
* NewDatabase
* @param driver string
* @return *DB
**/
func NewDatabase(driver string) (*DB, error) {
	if driver == "" {
		return nil, logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	now := time.Now()
	return &DB{
		CreatedAt:   now,
		UpdateAt:    now,
		Id:          utility.UUID(),
		Description: "",
		Schemas:     map[string]*Schema{},
		driver:      drivers[driver](),
	}, nil
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
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Connect(params)
}

/**
* Disconected
* @return error
**/
func (s *DB) Disconected() error {
	if s.driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
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
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.SetMain(params)
}

/**
* SetAdmin
* @return error
**/
func (s *DB) SetUser(username, password, confirmation string) error {
	if s.driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.SetUser(username, password, confirmation)
}

/**
* DeleteUser
* @return error
**/
func (s *DB) DeleteUser(username string) error {
	if s.driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
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
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
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
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.DropSchema(name)
}

/**
* CreateModel
* @param model *Model
* @return error
**/
func (s *DB) CreateModel(model *Model) error {
	if s.driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.CreateModel(model)
}

/**
* MutateModel
* @param model *Model
* @return error
**/
func (s *DB) MutateModel(model *Model) error {
	if s.driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.MutateModel(model)
}

/**
* Exec
* @param sql string
* @param params ...interface{}
* @return error
**/
func (s *DB) Exec(sql string, params ...interface{}) error {
	if s.driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Exec(sql, params...)
}

/**
* SQL
* @param sql string
* @param params ...interface{}
* @return et.Items, error
**/
func (s *DB) SQL(sql string, params ...interface{}) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.SQL(sql, params...)
}

/**
* One
* @param sql string
* @param params ...interface{}
* @return et.Item, error
**/
func (s *DB) One(sql string, params ...interface{}) (et.Item, error) {
	if s.driver == nil {
		return et.Item{}, logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.One(sql, params...)
}

/**
* Query
* @param linq *Linq
* @return et.Items, error
**/
func (s *DB) Query(linq *Linq) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Query(linq)
}

/**
* Count
* @param linq *Linq
* @return int, error
**/
func (s *DB) Count(linq *Linq) (int, error) {
	if s.driver == nil {
		return 0, logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Count(linq)
}

/**
* Last
* @param linq *Linq
* @return et.Items, error
**/
func (s *DB) Last(linq *Linq) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Last(linq)
}

/**
* Current
* @param command *Command
* @return et.Items, error
**/
func (s *DB) Current(command *Command) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Current(command)
}

/**
* Command
* @param command *Command
* @return et.Item, error
**/
func (s *DB) Command(command *Command) (et.Item, error) {
	if s.driver == nil {
		return et.Item{}, logs.NewError(MSG_DRIVER_NOT_DEFINED)
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

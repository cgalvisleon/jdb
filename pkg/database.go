package jdb

import (
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/utility"
)

type Database struct {
	CreatedAt   time.Time          `json:"created_date"`
	UpdateAt    time.Time          `json:"update_date"`
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Schemas     map[string]*Schema `json:"schemas"`
	Driver      *Driver            `json:"driver"`
}

/**
* NewDatabase
* @return *Database
**/
func NewDatabase(name, driver string) (*Database, error) {
	if driver == "" {
		return nil, logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	now := time.Now()

	return &Database{
		CreatedAt:   now,
		UpdateAt:    now,
		Id:          utility.UUID(),
		Name:        name,
		Description: "",
		Schemas:     map[string]*Schema{},
		Driver:      Drivers[driver],
	}, nil
}

/**
* Describe
* @return et.Json
**/
func (s *Database) Describe() et.Json {
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
func (s *Database) Conected(params et.Json) error {
	if s.Driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return (*s.Driver).Connect(params)
}

/**
* Disconected
* @return error
**/
func (s *Database) Disconected() error {
	if s.Driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return (*s.Driver).Disconnect()
}

/**
* SetMain
* @param params et.Json
* @return error
**/
func (s *Database) SetMain(params et.Json) error {
	if s.Driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return (*s.Driver).SetMain(params)
}

/**
* SetAdmin
* @return error
**/
func (s *Database) SetUser(username, password, confirmation string) error {
	if s.Driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return (*s.Driver).SetUser(username, password, confirmation)
}

/**
* DeleteUser
* @return error
**/
func (s *Database) DeleteUser(username string) error {
	if s.Driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return (*s.Driver).DeleteUser(username)
}

/**
* SetSerie
* @return int64
**/
func (s *Database) SetSerie(tag string, val int) int64 {
	if s.Driver == nil {
		return 0
	}

	return (*s.Driver).SetSerie(tag, val)
}

/**
* GetSerie
* @return int64
**/
func (s *Database) GetSerie(tag string) int64 {
	if s.Driver == nil {
		return 0
	}

	return (*s.Driver).GetSerie(tag)
}

/**
* NextCode
* @param tag string
* @param format string "%08v" "PREFIX-%08v-SUFFIX"
* @return string
**/
func (s *Database) NextCode(tag, format string) string {
	if s.Driver == nil {
		return ""
	}

	return (*s.Driver).NextCode(tag, format)
}

/**
* CreateModel
* @param model *Model
* @return error
**/
func (s *Database) CreateModel(model *Model) error {
	if s.Driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return (*s.Driver).CreateModel(model)
}

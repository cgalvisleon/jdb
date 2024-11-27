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
func NewDatabase(name string) *Database {
	now := time.Now()

	return &Database{
		CreatedAt:   now,
		UpdateAt:    now,
		Id:          utility.UUID(),
		Name:        name,
		Description: "",
		Schemas:     map[string]*Schema{},
	}
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
func (s *Database) Conected(params et.Json) bool {
	return (*s.Driver).Connect(params) == nil
}

/**
* Init
* @return error
**/
func (s *Database) Init() error {
	return nil
}

/**
* SetAdmin
* @return error
**/
func (s *Database) SetUser(username, password, confirmation string) error {
	return (*s.Driver).SetUser(username, password, confirmation)
}

/**
* DeleteUser
* @return error
**/
func (s *Database) DeleteUser(username string) error {
	return (*s.Driver).DeleteUser(username)
}

/**
* SetMain
**/
func (s *Database) SetMain(db *Database) error {
	if s.Driver == nil {
		return logs.NewError(MSG_DRIVER_NOT_DEFINED)
	}

	return nil
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

package jdb

import (
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
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
	Main        *Database
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
func (s *Database) SetMain(db *Database) {
	s.Main = db
}

/**
* SetSerie
* @return int64
**/
func (s *Database) SetSerie(tag string, val int) int64 {
	if s.Main == nil {
		return (*s.Driver).SetSerie(tag, val)
	}

	return s.Main.SetSerie(tag, val)
}

/**
* GetSerie
* @return int64
**/
func (s *Database) GetSerie(tag string) int64 {
	if s.Main == nil {
		return (*s.Driver).GetSerie(tag)
	}

	return s.Main.GetSerie(tag)
}

/**
* GetCode
* @param tag string
* @param format string "%08v" "PREFIX-%08v-SUFFIX"
* @return string
**/
func (s *Database) GetCode(tag, format string) string {
	val := s.GetSerie(tag)

	if len(format) == 0 {
		return strs.Format("%08v", val)
	}

	return strs.Format(format, val)
}

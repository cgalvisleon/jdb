package jdb

import (
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

type Database struct {
	CreatedAt      time.Time          `json:"created_date"`
	UpdateAt       time.Time          `json:"update_date"`
	Id             string             `json:"id"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	Schemas        map[string]*Schema `json:"schemas"`
	Driver         *Driver            `json:"driver"`
	IndexField     string             `json:"index_field"`
	DataField      string             `json:"data_field"`
	ProjectField   string             `json:"project_field"`
	CreatedAtField string             `json:"created_at_field"`
	UpdateAtField  string             `json:"update_at_field"`
	StateField     string             `json:"state_field"`
	KeyField       string             `json:"key_field"`
	Main           *Database
}

/**
* NewDatabase
* @return *Database
**/
func NewDatabase(name, description string) *Database {
	now := time.Now()

	return &Database{
		CreatedAt:      now,
		UpdateAt:       now,
		Id:             utility.UUID(),
		Name:           name,
		Description:    description,
		Schemas:        map[string]*Schema{},
		IndexField:     "_id",
		DataField:      "data",
		ProjectField:   "project_id",
		CreatedAtField: "created_at",
		UpdateAtField:  "update_at",
		StateField:     "_state",
		KeyField:       "_idt",
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
* SetParams
* @return error
**/
func (s *Database) SetParams(data et.Json) error {
	return (*s.Driver).SetParams(data)
}

/**
* SetMain
**/
func (s *Database) SetMain(db *Database) {
	s.Main = db
}

/**
* SetIndex
* @return int64
**/
func (s *Database) SetIndex(tag string, val int) int64 {
	if s.Main == nil {
		return (*s.Driver).SetIndex(tag, val)
	}

	return s.Main.SetIndex(tag, val)
}

/**
* GetIndex
* @return int64
**/
func (s *Database) GetIndex(tag string) int64 {
	if s.Main == nil {
		return (*s.Driver).GetIndex(tag)
	}

	return s.Main.GetIndex(tag)
}

/**
* GetCode
* @param tag string
* @param format string "%08v" "PREFIX-%08v-SUFFIX"
* @return string
**/
func (s *Database) GetCode(tag, format string) string {
	val := s.GetIndex(tag)

	if len(format) == 0 {
		return strs.Format("%08v", val)
	}

	return strs.Format(format, val)
}

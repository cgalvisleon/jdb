package jdb

import (
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

func schemaName(name string) string {
	return strs.Lowcase(name)
}

type Schema struct {
	Db          *Database
	CreatedAt   time.Time         `json:"created_date"`
	UpdateAt    time.Time         `json:"update_date"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Models      map[string]*Model `json:"models"`
}

func NewSchema(database *Database, name, description string) *Schema {
	now := time.Now()

	result := &Schema{
		Db:          database,
		CreatedAt:   now,
		UpdateAt:    now,
		Name:        name,
		Description: description,
		Models:      map[string]*Model{},
	}

	database.Schemas[name] = result

	return result
}

/**
* Describe
* @return et.Json
**/
func (s *Schema) Describe() et.Json {
	result, err := et.Object(s)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* Init
* @return error
**/
func (s *Schema) Init() error {
	return nil
}

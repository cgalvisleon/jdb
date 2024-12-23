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
	Db          *DB               `json:"-"`
	CreatedAt   time.Time         `json:"created_date"`
	UpdateAt    time.Time         `json:"update_date"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Models      map[string]*Model `json:"models"`
}

func NewSchema(db *DB, name string) (*Schema, error) {
	now := time.Now()

	result := &Schema{
		Db:          db,
		CreatedAt:   now,
		UpdateAt:    now,
		Name:        strs.Lowcase(name),
		Description: "",
		Models:      map[string]*Model{},
	}
	err := result.Init()
	if err != nil {
		return nil, err
	}

	db.Schemas[name] = result

	return result, nil
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
	return s.Db.CreateSchema(s.Name)
}

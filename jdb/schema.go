package jdb

import (
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

type Schema struct {
	Db          *DB               `json:"-"`
	CreatedAt   time.Time         `json:"created_date"`
	UpdateAt    time.Time         `json:"update_date"`
	Id          string            `json:"id"`
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
		Id:          utility.RecordId("schema", ""),
		Name:        Name(name),
		Description: "",
		Models:      map[string]*Model{},
	}
	err := result.Init()
	if err != nil {
		return nil, err
	}

	db.Schemas[name] = result
	Jdb.Schemas[name] = result

	return result, nil
}

/**
* Describe
* @return et.Json
**/
func (s *Schema) Describe() et.Json {
	var models = make([]et.Json, 0)
	for _, model := range s.Models {
		models = append(models, model.Describe())
	}

	result := et.Json{
		"created_date": s.CreatedAt,
		"update_date":  s.UpdateAt,
		"id":           s.Id,
		"name":         s.Name,
		"description":  s.Description,
		"models":       models,
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

/**
* Low
* @return string
**/
func (s *Schema) Low() string {
	return strs.Lowcase(s.Name)
}

/**
* Up
* @return string
**/
func (s *Schema) Up() string {
	return strs.Uppcase(s.Name)
}

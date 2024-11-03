package jdb

import (
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

func TableName(schema, name string) string {
	return strs.Format(`%s.%s`, strs.Lowcase(schema), strs.Uppcase(name))
}

type Model struct {
	Db           *Database
	Schema       *Schema
	CreatedAt    time.Time              `json:"created_date"`
	UpdateAt     time.Time              `json:"update_date"`
	Name         string                 `json:"name"`
	Table        string                 `json:"table"`
	Description  string                 `json:"description"`
	Columns      map[string]*Column     `json:"columns"`
	Indices      map[string]*Index      `json:"indices"`
	Keys         map[string]*Key        `json:"keys"`
	Relations    map[string]*Relation   `json:"relations"`
	Dictionaries map[string]*Dictionary `json:"dictionaries"`
	SourceField  *Column                `json:"data_field"`
	BeforeInsert Trigger                `json:"before_insert"`
	AfterInsert  Trigger                `json:"after_insert"`
	BeforeUpdate Trigger                `json:"before_update"`
	AfterUpdate  Trigger                `json:"after_update"`
	BeforeDelete Trigger                `json:"before_delete"`
	AfterDelete  Trigger                `json:"after_delete"`
	Integrity    bool                   `json:"integrity"`
}

func NewModel(schema *Schema, name, description string) *Model {
	now := time.Now()

	result := &Model{
		Db:           schema.Db,
		Schema:       schema,
		CreatedAt:    now,
		UpdateAt:     now,
		Name:         name,
		Description:  description,
		Table:        TableName(schema.Name, name),
		Columns:      map[string]*Column{},
		Indices:      map[string]*Index{},
		Keys:         map[string]*Key{},
		Relations:    map[string]*Relation{},
		Dictionaries: map[string]*Dictionary{},
	}

	result.Schema = schema

	return result
}

/**
* Describe
* @return et.Json
**/
func (s *Model) Describe() et.Json {
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
func (s *Model) Init() error {
	return nil
}

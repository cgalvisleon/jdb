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
	Details      map[string]*Column     `json:"details"`
	Indices      map[string]*Index      `json:"indices"`
	Uniques      map[string]*Index      `json:"uniques"`
	Keys         map[string]*Column     `json:"keys"`
	Relations    map[string]*Relation   `json:"relations"`
	Dictionaries map[string]*Dictionary `json:"dictionaries"`
	SourceField  *Column                `json:"data_field"`
	BeforeInsert []Trigger              `json:"before_insert"`
	AfterInsert  []Trigger              `json:"after_insert"`
	BeforeUpdate []Trigger              `json:"before_update"`
	AfterUpdate  []Trigger              `json:"after_update"`
	BeforeDelete []Trigger              `json:"before_delete"`
	AfterDelete  []Trigger              `json:"after_delete"`
	Functions    map[string]*Function   `json:"functions"`
	Integrity    bool                   `json:"integrity"`
}

func NewModel(schema *Schema, name string) *Model {
	now := time.Now()

	result := &Model{
		Db:           schema.Db,
		Schema:       schema,
		CreatedAt:    now,
		UpdateAt:     now,
		Name:         strs.Titlecase(name),
		Description:  "",
		Table:        TableName(schema.Name, name),
		Columns:      make(map[string]*Column),
		Details:      make(map[string]*Column),
		Indices:      make(map[string]*Index),
		Uniques:      make(map[string]*Index),
		Keys:         make(map[string]*Column),
		Relations:    make(map[string]*Relation),
		Dictionaries: make(map[string]*Dictionary),
		SourceField:  nil,
		BeforeInsert: []Trigger{},
		AfterInsert:  []Trigger{},
		BeforeUpdate: []Trigger{},
		AfterUpdate:  []Trigger{},
		BeforeDelete: []Trigger{},
		AfterDelete:  []Trigger{},
		Functions:    make(map[string]*Function),
		Integrity:    false,
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

/**
* DefineColumn
* @param name string
* @return *Column
**/
func (s *Model) GetColumn(name string) *Column {
	field := fieldName(name)
	if col, ok := s.Columns[field]; ok {
		return col
	}

	return nil
}

/**
* DefineColumn
* @param name string
* @return *Column
**/
func (s *Model) GetColumns(names ...string) []*Column {
	result := []*Column{}
	for _, name := range names {
		if col := s.GetColumn(name); col != nil {
			result = append(result, col)
		}
	}

	return result
}

func (s *Model) ExecDetails(data *et.Json) error {

	return nil
}

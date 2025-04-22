package jdb

import (
	"encoding/json"
	"slices"
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/strs"
)

type Schema struct {
	Db          *DB       `json:"-"`
	CreatedAt   time.Time `json:"created_date"`
	UpdateAt    time.Time `json:"update_date"`
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UseCore     bool      `json:"use_core"`
	models      []*Model  `json:"-"`
	flows       []*Flow   `json:"-"`
}

/**
* NewSchema
* @param db *DB, name string
* @return *Schema, error
**/
func NewSchema(db *DB, name string) (*Schema, error) {
	name = Name(name)
	idx := slices.IndexFunc(db.schemas, func(schema *Schema) bool { return schema.Name == name })
	if idx != -1 {
		return db.schemas[idx], nil
	}

	now := time.Now()
	result := &Schema{
		Db:        db,
		CreatedAt: now,
		UpdateAt:  now,
		Id:        reg.Id("schema"),
		Name:      name,
		UseCore:   db.UseCore,
		models:    make([]*Model, 0),
		flows:     make([]*Flow, 0),
	}
	err := result.Init()
	if err != nil {
		return nil, err
	}

	db.schemas = append(db.schemas, result)

	return result, nil
}

/**
* Save
* @return error
**/
func (s *Schema) Save() error {
	if !s.UseCore {
		return mistake.New(MSG_SCHEMA_NOT_USING_CORE)
	}

	buf, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = s.Db.upsertModel("schema", s.Name, 1, buf)
	if err != nil {
		return err
	}

	return nil
}

/**
* Describe
* @return et.Json
**/
func (s *Schema) Describe() et.Json {
	var models = make([]et.Json, 0)
	for _, model := range s.models {
		models = append(models, model.Describe())
	}
	var flows = make([]et.Json, 0)
	for _, flow := range s.flows {
		flows = append(flows, flow.Describe())
	}

	result := et.Json{
		"created_date": s.CreatedAt,
		"update_date":  s.UpdateAt,
		"id":           s.Id,
		"name":         s.Name,
		"description":  s.Description,
		"models":       models,
		"flows":        flows,
	}

	return result
}

/**
* GetModelByProjectId
* @param name, projectId string
* @return *Model
**/
func (s *Schema) GetModelByProjectId(name, projectId string) *Model {
	name = Name(name)
	idx := slices.IndexFunc(s.models, func(model *Model) bool { return model.Name == name && model.projectId == projectId })
	if idx != -1 {
		return s.models[idx]
	}

	idx = slices.IndexFunc(s.models, func(model *Model) bool { return model.Name == name })
	if idx != -1 {
		return s.models[idx]
	}

	return nil
}

/**
* GetModel
* @param name string
* @return *Model
**/
func (s *Schema) GetModel(name string) *Model {
	return s.GetModelByProjectId(name, "")
}

/**
* GetFlow
* @param name string
* @return *Flow
**/
func (s *Schema) GetFlow(name string) *Flow {
	idx := slices.IndexFunc(s.flows, func(flow *Flow) bool { return flow.Name == name })
	if idx != -1 {
		return s.flows[idx]
	}
	return nil
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

package jdb

import (
	"encoding/json"
	"slices"
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/timezone"
	"github.com/cgalvisleon/et/utility"
)

type Schema struct {
	Db          *DB       `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"update_at"`
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UseCore     bool      `json:"use_core"`
	models      []*Model  `json:"-"`
}

/**
* NewSchema
* @param db *DB, name string
* @return *Schema
**/
func NewSchema(db *DB, name string) *Schema {
	idx := slices.IndexFunc(db.schemas, func(e *Schema) bool { return e.Name == name })
	if idx != -1 {
		return db.schemas[idx]
	}

	now := timezone.NowTime()
	result := &Schema{
		Db:        db,
		CreatedAt: now,
		UpdateAt:  now,
		Id:        utility.UUID(),
		Name:      name,
		UseCore:   db.UseCore,
		models:    make([]*Model, 0),
	}
	db.addSchema(result)

	return result
}

/**
* AddModel
* @param model *Model
**/
func (s *Schema) addModel(model *Model) {
	idx := slices.IndexFunc(s.Db.models, func(e *Model) bool { return e.Name == model.Name })
	if idx == -1 {
		s.Db.models = append(s.Db.models, model)
	}

	idx = slices.IndexFunc(s.models, func(e *Model) bool { return e.Name == model.Name })
	if idx == -1 {
		s.models = append(s.models, model)
	}
}

/**
* DropModel
* @param model *Model
**/
func (s *Schema) dropModel(model *Model) {
	idx := slices.IndexFunc(s.Db.models, func(e *Model) bool { return e.Name == model.Name })
	if idx != -1 {
		s.Db.models = append(s.Db.models[:idx], s.Db.models[idx+1:]...)
	}

	idx = slices.IndexFunc(s.models, func(e *Model) bool { return e.Name == model.Name })
	if idx != -1 {
		s.models = append(s.models[:idx], s.models[idx+1:]...)
	}
}

/**
* Serialize
* @return []byte, error
**/
func (s *Schema) serialize() ([]byte, error) {
	result, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

/**
* Describe
* @return et.Json
**/
func (s *Schema) Describe() et.Json {
	definition, err := s.serialize()
	if err != nil {
		return et.Json{}
	}

	result := et.Json{}
	err = json.Unmarshal(definition, &result)
	if err != nil {
		return et.Json{}
	}

	var models = make([]et.Json, 0)
	for _, model := range s.models {
		models = append(models, model.Describe())
	}

	result["kind"] = "schema"
	result["models"] = models

	return result
}

/**
* Mutate
* @return error
**/
func (s *Schema) Drop() error {
	return s.Db.DropSchema(s.Name)
}

/**
* GetModel
* @param name string
* @return *Model
**/
func (s *Schema) GetModel(name string) *Model {
	idx := slices.IndexFunc(s.models, func(e *Model) bool { return e.Name == name })
	if idx != -1 {
		return s.models[idx]
	}

	return nil
}

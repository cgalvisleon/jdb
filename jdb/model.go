package jdb

import (
	"encoding/gob"
	"encoding/json"
	"slices"
	"time"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

func TableName(schema, name string) string {
	return strs.Format(`%s.%s`, strs.Lowcase(schema), strs.Lowcase(name))
}

type Model struct {
	Db             *DB                    `json:"-"`
	Schema         *Schema                `json:"-"`
	CreatedAt      time.Time              `json:"created_date"`
	UpdateAt       time.Time              `json:"update_date"`
	Name           string                 `json:"name"`
	Table          string                 `json:"table"`
	Description    string                 `json:"description"`
	Columns        []*Column              `json:"columns"`
	Indices        map[string]*Index      `json:"indices"`
	Uniques        map[string]*Index      `json:"uniques"`
	Keys           map[string]*Column     `json:"keys"`
	References     []*Reference           `json:"references"`
	Dictionaries   map[string]*Dictionary `json:"dictionaries"`
	ColRequired    map[string]bool        `json:"col_required"`
	SourceField    *Column                `json:"data_field"`
	SystemKeyField *Column                `json:"key_field"`
	StateField     *Column                `json:"state_field"`
	IndexField     *Column                `json:"index_field"`
	ClassField     *Column                `json:"class_field"`
	FullTextField  *Column                `json:"full_text"`
	BeforeInsert   []string               `json:"-"`
	AfterInsert    []string               `json:"-"`
	BeforeUpdate   []string               `json:"-"`
	AfterUpdate    []string               `json:"-"`
	BeforeDelete   []string               `json:"-"`
	AfterDelete    []string               `json:"-"`
	Functions      map[string]*Function   `json:"-"`
	Integrity      bool                   `json:"integrity"`
	Version        int                    `json:"version"`
}

/**
* NewModel
* @param schema *Schema
* @param name string
* @return *Model
**/
func NewModel(schema *Schema, name string, version int) *Model {
	if version == 0 {
		version = 1
	}
	now := time.Now()
	name = Name(name)
	result := &Model{
		Db:           schema.Db,
		Schema:       schema,
		CreatedAt:    now,
		UpdateAt:     now,
		Name:         name,
		Description:  "",
		Table:        TableName(schema.Name, name),
		Columns:      make([]*Column, 0),
		Indices:      make(map[string]*Index),
		Uniques:      make(map[string]*Index),
		Keys:         make(map[string]*Column),
		References:   make([]*Reference, 0),
		Dictionaries: make(map[string]*Dictionary),
		BeforeInsert: []string{},
		AfterInsert:  []string{},
		BeforeUpdate: []string{},
		AfterUpdate:  []string{},
		BeforeDelete: []string{},
		AfterDelete:  []string{},
		Functions:    make(map[string]*Function),
		Integrity:    false,
		Version:      version,
	}

	schema.Models[result.Name] = result

	return result
}

func init() {
	gob.Register(&Column{})
}

/**
* Up
* @return string
**/
func (s *Model) Up() string {
	return strs.Uppcase(s.Name)
}

/**
* Low
* @return string
**/
func (s *Model) Low() string {
	return strs.Lowcase(s.Name)
}

/**
* Serialized
* @return []byte, error
**/
func (s *Model) Serialized() ([]byte, error) {
	return json.Marshal(s)
}

/**
* Deserialize
* @param data []byte
* @return error
**/
func (s *Model) Load(data []byte) error {
	err := json.Unmarshal(data, s)
	if err != nil {
		return err
	}

	return nil
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
	if s.Db == nil {
		return console.Alertm(MSG_DATABASE_IS_REQUIRED)
	}

	return s.Db.LoadModel(s)
}

/**
* DefineColumn
* @param name string
* @return *Column
**/
func (s *Model) GetColumn(name string) *Column {
	for _, col := range s.Columns {
		if col.Low() == strs.Lowcase(name) {
			return col
		}
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

/**
* GetKeys
* @return []*Column
**/
func (s *Model) GetKeys() []*Column {
	result := []*Column{}
	for _, col := range s.Keys {
		result = append(result, col)
	}

	return result
}

func (s *Model) GetForeignKeys() []string {

	return nil
}

/**
* GetDetails
* @return []et.Json
**/
func (s *Model) GetDetails(data *et.Json) *et.Json {
	for _, col := range s.Columns {
		switch col.TypeColumn {
		case TpGenerate:
			col.Definition.(FuncGenerated)(col, data)
		case TpDetail:
			linq := From(col.Definition.(*Model))
			first := true
			for _, key := range col.Model.Keys {
				val := (*data)[key.Field]
				if first {
					linq = linq.Where(key.Fk()).Eq(val)
					first = false
				} else {
					linq = linq.And(key.Fk()).Eq(val)
				}
			}
			result, err := linq.Page(1).Rows(30)
			if err == nil {
				data.Set(col.Low(), result)
			}
		}
	}
	return data
}

/**
* Column
* @param name string
* @return *Column
**/
func (s *Model) Column(name string) *Column {
	return s.GetColumn(name)
}

/**
* Col
* @param name string
* @return *Column
**/
func (s *Model) Col(name string) *Column {
	return s.GetColumn(name)
}

/**
* New
* @param data et.Json
* @return et.Json
**/
func (s *Model) New() et.Json {
	var result = &et.Json{}
	var details = []*Column{}
	for _, col := range s.Columns {
		if slices.Contains([]*Column{s.SystemKeyField, s.IndexField}, col) {
			continue
		}
		switch col.TypeColumn {
		case TpGenerate:
			col.Definition.(FuncGenerated)(col, result)
		case TpDetail:
			details = append(details, col)
		default:
			result.Set(col.Low(), col.DefaultValue())
		}
	}

	for _, col := range details {
		dtl := col.Definition.(*Model).New()
		for _, key := range col.Model.Keys {
			val := (*result)[key.Field]
			dtl.Set(key.Fk(), val)
		}
		result.Set(col.Low(), dtl)
	}

	return *result
}

/**
* Col
* @param name string
* @return *Column
**/
func (s *Model) FullText() *Column {
	return s.FullTextField
}

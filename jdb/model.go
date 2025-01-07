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
	SourceField    *Column                `json:"source_field"`
	SystemKeyField *Column                `json:"system_key_field"`
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
		ColRequired:  make(map[string]bool),
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
	models[result.Table] = result

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
* GetColumn
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
* GetColumns
* @param name string
* @return *Column
*
 */
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
* SetField
* @param name string
* @return *Field
**/
func (s *Model) SetField(name string, create bool) *Field {
	col := s.GetColumn(name)
	if col == nil && !create {
		return nil
	}

	if col == nil {
		col = s.DefineAtribute(name, TypeDataText)
	}
	result := &Field{
		Column: col,
		Schema: s.Schema.Name,
		Table:  s.Name,
		Field:  col.Field,
		Name:   col.Name,
		Alias:  col.Name,
	}

	if col.TypeColumn != TpColumn {
		result.Atrib = col.Name
	}

	return result
}

/**
* GetField
* @param name string
* @return *Field
**/
func (s *Model) GetField(name string) *Field {
	list := strs.Split(name, ".")
	switch len(list) {
	case 1:
		return s.SetField(list[0], false)
	case 2:
		if s.Name != strs.Lowcase(list[0]) {
			return nil
		}
		return s.SetField(list[1], !s.Integrity)
	case 3:
		table := strs.Format(`%s.%s`, list[0], list[1])
		if s.Table != strs.Lowcase(table) {
			return nil
		}
		return s.SetField(list[2], !s.Integrity)
	default:
		return nil
	}
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
			model := col.Definition.(*Model)
			var filter FilterTo
			linq := From(model)
			for _, key := range col.Model.Keys {
				val := (*data)[key.Field]
				if filter == nil {
					filter = linq.Where(key.Fk()).Eq(val)
				} else {
					filter = linq.And(key.Fk()).Eq(val)
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

/**
* MakeCollection
* @return *Model
**/
func (s *Model) MakeCollection() *Model {
	s.DefineColumn(CreatedAtField.Str(), CreatedAtField.TypeData())
	s.DefineColumn(UpdatedAtField.Str(), UpdatedAtField.TypeData())
	s.DefineColumn(ProjectField.Str(), ProjectField.TypeData())
	s.DefineColumn(StateField.Str(), StateField.TypeData())
	s.DefineColumn(KeyField.Str(), KeyField.TypeData())
	class := s.DefineColumn(ClassField.Str(), ClassField.TypeData())
	class.Default = s.Low()
	s.DefineColumn(SourceField.Str(), SourceField.TypeData())
	s.DefineColumn(SystemKeyField.Str(), SystemKeyField.TypeData())
	s.DefineColumn(IndexField.Str(), IndexField.TypeData())
	s.DefineKey(KeyField.Str())

	return s
}

/**
* MakeDetail
* @param fkeys []*Column
* @return *Model
**/
func (s *Model) MakeDetail(fkeys []*Column) *Model {
	s.DefineColumn(CreatedAtField.Str(), CreatedAtField.TypeData())
	s.DefineColumn(UpdatedAtField.Str(), UpdatedAtField.TypeData())
	s.DefineColumn(KeyField.Str(), KeyField.TypeData())
	s.DefineColumn(SourceField.Str(), SourceField.TypeData())
	s.DefineColumn(SystemKeyField.Str(), SystemKeyField.TypeData())
	s.DefineColumn(IndexField.Str(), IndexField.TypeData())
	s.DefineKey(KeyField.Str())
	for _, key := range fkeys {
		fkn := key.Fk()
		fk := s.DefineColumn(fkn, key.TypeData)
		NewReference(fk, RelationManyToOne, key)
		ref := NewReference(key, RelationOneToMany, fk)
		ref.OnDeleteCascade = true
		ref.OnUpdateCascade = true
		s.DefineRequired(fkn)
	}

	return s
}

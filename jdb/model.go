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
	Dictionaries   map[string]*Dictionary `json:"-"`
	ColRequired    map[string]bool        `json:"col_required"`
	KeyField       *Column                `json:"key_field"`
	SourceField    *Column                `json:"source_field"`
	SystemKeyField *Column                `json:"system_key_field"`
	StateField     *Column                `json:"state_field"`
	IndexField     *Column                `json:"index_field"`
	ClassField     *Column                `json:"class_field"`
	FullText       []string               `json:"full_text"`
	BeforeInsert   []string               `json:"-"`
	AfterInsert    []string               `json:"-"`
	BeforeUpdate   []string               `json:"-"`
	AfterUpdate    []string               `json:"-"`
	BeforeDelete   []string               `json:"-"`
	AfterDelete    []string               `json:"-"`
	Details        map[string]*Model      `json:"-"`
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
	table := TableName(schema.Name, name)
	result := models[table]
	if result != nil {
		return result
	}

	result = &Model{
		Db:           schema.Db,
		Schema:       schema,
		CreatedAt:    now,
		UpdateAt:     now,
		Name:         name,
		Description:  "",
		Table:        table,
		Columns:      make([]*Column, 0),
		Indices:      make(map[string]*Index),
		Uniques:      make(map[string]*Index),
		Keys:         make(map[string]*Column),
		References:   make([]*Reference, 0),
		Dictionaries: make(map[string]*Dictionary),
		ColRequired:  make(map[string]bool),
		FullText:     []string{},
		BeforeInsert: []string{},
		AfterInsert:  []string{},
		BeforeUpdate: []string{},
		AfterUpdate:  []string{},
		BeforeDelete: []string{},
		AfterDelete:  []string{},
		Details:      make(map[string]*Model),
		Functions:    make(map[string]*Function),
		Integrity:    false,
		Version:      version,
	}

	schema.Models[result.Name] = result
	models[table] = result

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

	dictionaries := map[string]et.Json{}
	for key, value := range s.Dictionaries {
		dictionaries[key] = value.Describe()
	}
	result["dictionaries"] = dictionaries

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
func (s *Model) SetField(name string, isCreated bool) *Field {
	col := s.GetColumn(name)
	if col == nil && !isCreated {
		return nil
	}

	if col == nil {
		col = s.DefineAtribute(name, TypeDataText)
	}

	return col.GetField()
}

/**
* GetField
* @param name string, isCreated bool
* @return *Field
**/
func (s *Model) GetField(name string, isCreated bool) *Field {
	list := strs.Split(name, ".")
	switch len(list) {
	case 1:
		isCreated = isCreated && s.SourceField != nil
		return s.SetField(list[0], isCreated)
	case 2:
		if s.Name != strs.Lowcase(list[0]) {
			return nil
		}
		isCreated = s.SourceField != nil && !s.Integrity
		return s.SetField(list[1], isCreated)
	case 3:
		table := strs.Format(`%s.%s`, list[0], list[1])
		if s.Table != strs.Lowcase(table) {
			return nil
		}
		isCreated = s.SourceField != nil && !s.Integrity
		return s.SetField(list[2], isCreated)
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
	if data == nil {
		data = &et.Json{}
	} else if data.IsEmpty() {
		data = &et.Json{}
	}

	for _, col := range s.Columns {
		switch col.TypeColumn {
		case TpGenerate:
			if col.FuncGenerated != nil {
				col.FuncGenerated(col, data)
			}
		case TpDetail:
			model := col.Detail
			if model == nil {
				continue
			}
			var filter FilterTo
			linq := From(model)
			for _, key := range col.Model.Keys {
				val := (*data)[key.Field]
				if val == nil {
					break
				}
				if filter == nil {
					filter = linq.Where(key.Fk()).Eq(val)
				} else {
					filter = linq.And(key.Fk()).Eq(val)
				}
			}
			result, err := linq.
				Page(1).
				Rows(col.Limit)
			if err != nil {
				data.Set(col.Name, result)
			} else {
				data.Set(col.Name, []et.Json{})
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
	for _, col := range s.Columns {
		if slices.Contains([]*Column{s.SystemKeyField, s.IndexField}, col) {
			continue
		}
		switch col.TypeColumn {
		case TpGenerate:
			if col.FuncGenerated != nil {
				col.FuncGenerated(col, result)
			}
		case TpColumn:
			if col != s.SourceField {
				result.Set(col.Name, col.DefaultValue())
			}
		default:
			result.Set(col.Name, col.DefaultValue())
		}
	}

	for _, detail := range s.Details {
		dtl := detail.New()
		for _, key := range detail.Keys {
			val := (*result)[key.Field]
			dtl.Set(key.Fk(), val)
		}
		result.Set(detail.Name, dtl)
	}

	return *result
}

/**
* MakeCollection
* @return *Model
**/
func (s *Model) MakeCollection() *Model {
	s.DefineCreatedAtField()
	s.DefineUpdatedAtField()
	s.DefineStateField()
	s.DefineKeyField()
	s.DefineClassField()
	s.DefineSourceField()
	s.DefineSystemKeyField()
	s.DefineIndexField()

	return s
}

/**
* MakeDetailRelation
* @param fkeys []*Column
* @return *Model
**/
func (s *Model) MakeDetailRelation(owner *Model) *Model {
	key := owner.KeyField
	fkn := owner.Name
	fk := s.DefineColumn(fkn, key.TypeData)
	NewReference(fk, RelationManyToOne, key)
	ref := NewReference(key, RelationOneToMany, fk)
	ref.OnDeleteCascade = true
	ref.OnUpdateCascade = true
	s.DefineRequired(fkn)
	s.DefineSystemKeyField()
	s.DefineIndexField()

	return s
}

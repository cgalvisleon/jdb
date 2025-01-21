package jdb

import (
	"encoding/gob"
	"encoding/json"
	"slices"
	"strings"
	"time"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

type RID struct {
	Id     string `json:"id"`
	Schema string `json:"schema"`
	Table  string `json:"table"`
	Model  string `json:"model"`
}

func GetRID(id string) *RID {
	result := &RID{
		Id: id,
	}

	split := strings.Split(id, ":")
	if len(split) == 2 {
		result.Table = split[0]
		split = strings.Split(split[0], ".")
		if len(split) == 2 {
			result.Schema = split[0]
			result.Model = split[1]
		}
	}

	return result
}

func TableName(schema, name string) string {
	return strs.Format(`%s.%s`, strs.Lowcase(schema), strs.Lowcase(name))
}

type Model struct {
	Db             *DB                         `json:"-"`
	Schema         *Schema                     `json:"-"`
	CreatedAt      time.Time                   `json:"created_date"`
	UpdateAt       time.Time                   `json:"update_date"`
	Name           string                      `json:"name"`
	Table          string                      `json:"table"`
	Description    string                      `json:"description"`
	Columns        []*Column                   `json:"columns"`
	Indices        map[string]*Index           `json:"indices"`
	Uniques        map[string]*Index           `json:"uniques"`
	Keys           map[string]*Column          `json:"keys"`
	References     []*Reference                `json:"references"`
	Dictionaries   map[interface{}]*Dictionary `json:"-"`
	ColRequired    map[string]bool             `json:"col_required"`
	CreatedAtField *Column                     `json:"created_at_field"`
	UpdatedAtField *Column                     `json:"updated_at_field"`
	KeyField       *Column                     `json:"key_field"`
	SourceField    *Column                     `json:"source_field"`
	SystemKeyField *Column                     `json:"system_key_field"`
	StateField     *Column                     `json:"state_field"`
	IndexField     *Column                     `json:"index_field"`
	FullTextField  *Column                     `json:"full_text_field"`
	EventsInsert   []Event                     `json:"-"`
	EventsUpdate   []Event                     `json:"-"`
	EventsDelete   []Event                     `json:"-"`
	Details        map[string]*Model           `json:"-"`
	Functions      map[string]*Function        `json:"-"`
	Integrity      bool                        `json:"integrity"`
	Version        int                         `json:"version"`
	Show           bool                        `json:"-"`
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
		Dictionaries: make(map[interface{}]*Dictionary),
		ColRequired:  make(map[string]bool),
		EventsInsert: make([]Event, 0),
		EventsUpdate: make([]Event, 0),
		EventsDelete: make([]Event, 0),
		Details:      make(map[string]*Model),
		Functions:    make(map[string]*Function),
		Integrity:    false,
		Version:      version,
	}
	if schema.Db.UseCore {
		result.DefineSystemKeyField()
		result.DefineIndexField()
	}
	schema.Models[result.Name] = result
	models[table] = result

	return result
}

func init() {
	gob.Register(&Column{})
}

/**
* GenId
* @param id string
* @return string
**/
func (s *Model) GenId(id string) string {
	if !map[string]bool{"": true, "*": true, "new": true}[id] {
		split := strings.Split(id, ":")
		if len(split) == 1 {
			return strs.Format(`%s:%s`, s.Table, id)
		}

		return id
	}

	id = utility.Snowflake()
	return strs.Format(`%s:%s`, s.Table, id)
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

	dictionaries := map[interface{}]et.Json{}
	for key, value := range s.Dictionaries {
		dictionaries[key] = value.Describe()
	}
	result["dictionaries"] = dictionaries

	return result
}

/**
* Debug
* @return *Model
**/
func (s *Model) Debug() *Model {
	s.Show = true

	return s
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
		if col.Name == name {
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
* New
* @param data et.Json
* @return et.Json
**/
func (s *Model) New(data et.Json) et.Json {
	var result = &et.Json{}
	defaultColValue := func(cols []*Column) {
		for _, col := range cols {
			if slices.Contains([]*Column{s.SystemKeyField, s.IndexField}, col) {
				continue
			}
			switch col.TypeColumn {
			case TpColumn:
				if col != s.SourceField {
					result.Set(col.Name, col.DefaultValue())
				}
			case TpAtribute:
				result.Set(col.Name, col.DefaultValue())
			case TpGenerated:
				if col.FuncGenerated != nil {
					col.FuncGenerated(col, result)
				}
			}
		}
	}

	defaultColValue(s.Columns)

	for _, value := range data {
		dictionary := s.Dictionaries[value]
		if dictionary != nil {
			defaultColValue(dictionary.Columns)
		}
	}

	for _, detail := range s.Details {
		dtl := detail.New(et.Json{})
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
	s.DefineSystemKeyField()
	s.DefineCreatedAtField()
	s.DefineUpdatedAtField()
	s.DefineStateField()
	s.DefineKeyField()
	s.DefineSourceField()
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

	return s
}

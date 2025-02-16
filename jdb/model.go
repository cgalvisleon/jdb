package jdb

import (
	"encoding/json"
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

func TableName(schema *Schema, name string) string {
	if schema != nil {
		return strs.Format(`%s.%s`, strs.Lowcase(schema.Name), strs.Lowcase(name))
	}

	return strs.Lowcase(name)
}

type Model struct {
	Db              *DB                  `json:"-"`
	Schema          *Schema              `json:"-"`
	CreatedAt       time.Time            `json:"created_date"`
	UpdateAt        time.Time            `json:"update_date"`
	Name            string               `json:"name"`
	Table           string               `json:"table"`
	Description     string               `json:"description"`
	Columns         []*Column            `json:"columns"`
	GeneratedFields []*Column            `json:"generated_fields"`
	PrimaryKeys     []*Column            `json:"primary_keys"`
	ForeignKeys     []*Column            `json:"foreign_keys"`
	Indices         []*Index             `json:"indices"`
	Uniques         []*Index             `json:"uniques"`
	Relations       map[string]*Relation `json:"Relations"`
	Details         map[string]*Relation `json:"details"`
	History         *Relation            `json:"history"`
	ColRequired     map[string]bool      `json:"col_required"`
	SystemKeyField  *Column              `json:"system_key_field"`
	StateField      *Column              `json:"state_field"`
	IndexField      *Column              `json:"index_field"`
	SourceField     *Column              `json:"source_field"`
	FullTextField   *Column              `json:"full_text_field"`
	EventsInsert    []Event              `json:"-"`
	EventsUpdate    []Event              `json:"-"`
	EventsDelete    []Event              `json:"-"`
	Integrity       bool                 `json:"integrity"`
	IsCreated       bool                 `json:"is_created"`
	Version         int                  `json:"version"`
	IsDebug         bool                 `json:"-"`
}

/**
* NewModel
* @param schema *Schema, name string, tp TypeModel, version int
* @return *Model
**/
func NewModel(schema *Schema, name string, version int) *Model {
	if version == 0 {
		version = 1
	}
	now := time.Now()
	name = Name(name)
	table := TableName(schema, name)
	result := Jdb.Models[table]
	if result != nil {
		result.Version = version

		return result
	}

	result = &Model{
		Db:              schema.Db,
		Schema:          schema,
		CreatedAt:       now,
		UpdateAt:        now,
		Name:            name,
		Table:           table,
		Description:     "",
		Columns:         make([]*Column, 0),
		GeneratedFields: make([]*Column, 0),
		PrimaryKeys:     make([]*Column, 0),
		ForeignKeys:     make([]*Column, 0),
		Indices:         make([]*Index, 0),
		Uniques:         make([]*Index, 0),
		Relations:       make(map[string]*Relation),
		Details:         make(map[string]*Relation),
		History:         &Relation{Limit: -1},
		ColRequired:     make(map[string]bool),
		EventsInsert:    make([]Event, 0),
		EventsUpdate:    make([]Event, 0),
		EventsDelete:    make([]Event, 0),
		Version:         version,
	}
	result.DefineEvent(EventInsert, EventInsertDefault)
	result.DefineEvent(EventUpdate, EventUpdateDefault)
	result.DefineEvent(EventDelete, EventDeleteDefault)
	schema.Models[result.Name] = result
	Jdb.Models[table] = result
	result.IsCreated, _ = result.Db.LoadTable(result)

	return result
}

/**
* GetFrom
* @return *QlFrom
**/
func (s *Model) GetFrom() *QlFrom {
	return &QlFrom{Model: s}
}

/**
* GenId
* @param id string
* @return string
**/
func (s *Model) GenId(id string) string {
	if !map[string]bool{"": true, "*": true, "new": true}[id] {
		return id
	}

	return utility.RecordId(s.Table, id)
}

/**
* GenKey
* @param id string
* @return string
**/
func (s *Model) GenKey(id string) string {
	return utility.RecordId(s.Table, id)
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
* Serialized
* @return []byte, error
**/
func (s *Model) Serialized() ([]byte, error) {
	obj := s.Describe()

	if s.IsDebug {
		console.Debug(obj.ToString())
	}

	return json.Marshal(obj)
}

/**
* GetSerie
* @return int
**/
func (s *Model) GetSerie() int64 {
	return s.Db.GetSerie(s.Table)
}

/**
* Load
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
* Debug
* @return *Model
**/
func (s *Model) Debug() *Model {
	s.IsDebug = true

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
* Drop
**/
func (s *Model) Drop() {
	if s.Db == nil {
		return
	}

	for _, detail := range s.Details {
		detail.With.Drop()
	}

	s.Db.DropModel(s)
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
* getField
* @param name string
* @return *Field
**/
func (s *Model) getField(name string) *Field {
	col := s.GetColumn(name)
	if col == nil && s.Integrity {
		return nil
	}

	if s.SourceField == nil {
		return nil
	}

	col = s.DefineAtribute(name, TypeDataText)

	return col.GetField()
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
		return s.getField(list[0])
	case 2:
		if s.Name != strs.Lowcase(list[0]) {
			return nil
		}
		return s.getField(list[1])
	case 3:
		table := strs.Format(`%s.%s`, list[0], list[1])
		if s.Table != strs.Lowcase(table) {
			return nil
		}
		return s.getField(list[2])
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
	result = append(result, s.PrimaryKeys...)

	return result
}

/**
* Where
* @param val interface{}
* @return *Ql
**/
func (s *Model) Where(val interface{}) *Ql {
	return From(s)
}

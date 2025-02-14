package jdb

import (
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

type TypeModel int

const (
	TpTable TypeModel = iota
	TpModel
	TpCollection
)

type Model struct {
	Db             *DB                      `json:"-"`
	Schema         *Schema                  `json:"-"`
	Type           TypeModel                `json:"type"`
	CreatedAt      time.Time                `json:"created_date"`
	UpdateAt       time.Time                `json:"update_date"`
	Name           string                   `json:"name"`
	Table          string                   `json:"table"`
	Description    string                   `json:"description"`
	Columns        []*Column                `json:"columns"`
	Indices        map[string]*Index        `json:"indices"`
	Uniques        map[string]*Index        `json:"uniques"`
	Keys           map[string]*Column       `json:"keys"`
	ForeignKeys    map[string]*Reference    `json:"foreign_keys"`
	References     []*Reference             `json:"references"`
	Relations      []*Relation              `json:"relations"`
	Dictionaries   map[string][]*Dictionary `json:"dictionaries"`
	ColRequired    map[string]bool          `json:"col_required"`
	CreatedAtField *Column                  `json:"created_at_field"`
	UpdatedAtField *Column                  `json:"updated_at_field"`
	KeyField       *Column                  `json:"key_field"`
	Source         string                   `json:"source"`
	SourceField    *Column                  `json:"source_field"`
	SystemKeyField *Column                  `json:"system_key_field"`
	StateField     *Column                  `json:"state_field"`
	IndexField     *Column                  `json:"index_field"`
	FullTextField  *Column                  `json:"full_text_field"`
	EventsInsert   []Event                  `json:"-"`
	EventsUpdate   []Event                  `json:"-"`
	EventsDelete   []Event                  `json:"-"`
	Details        map[string]*Detail       `json:"-"`
	Integrity      bool                     `json:"integrity"`
	Log            int64                    `json:"log"`
	History        *Model                   `json:"-"`
	HistoryLimit   int64                    `json:"history_limit"`
	Version        int                      `json:"version"`
	Show           bool                     `json:"-"`
}

/**
* newModel
* @param schema *Schema, name string, tp TypeModel, version int
* @return *Model
**/
func newModel(schema *Schema, name string, tp TypeModel, version int) *Model {
	if version == 0 {
		version = 1
	}
	now := time.Now()
	name = Name(name)
	table := TableName(schema.Name, name)
	result := Jdb.Models[table]
	if result != nil {
		return result
	}

	result = &Model{
		Db:           schema.Db,
		Schema:       schema,
		Type:         tp,
		CreatedAt:    now,
		UpdateAt:     now,
		Name:         name,
		Description:  "",
		Table:        table,
		Columns:      make([]*Column, 0),
		Indices:      make(map[string]*Index),
		Uniques:      make(map[string]*Index),
		Keys:         make(map[string]*Column),
		ForeignKeys:  make(map[string]*Reference),
		References:   make([]*Reference, 0),
		Relations:    make([]*Relation, 0),
		Dictionaries: make(map[string][]*Dictionary),
		ColRequired:  make(map[string]bool),
		EventsInsert: make([]Event, 0),
		EventsUpdate: make([]Event, 0),
		EventsDelete: make([]Event, 0),
		Details:      make(map[string]*Detail),
		Source:       SOURCE,
		Integrity:    false,
		HistoryLimit: 0,
		Version:      version,
	}
	result.DefineEvent(EventInsert, EventInsertDefault)
	result.DefineEvent(EventUpdate, EventUpdateDefault)
	result.DefineEvent(EventDelete, EventDeleteDefault)
	if slices.Contains([]TypeModel{TpModel, TpCollection}, tp) {
		result.DefineIndexField()
		result.DefineSystemKeyField()
	}
	schema.Models[result.Name] = result
	Jdb.Models[table] = result

	return result
}

/**
* NewTable
* @param schema *Schema, name string, version int
* @return *Model
**/
func NewTable(schema *Schema, name string, version int) *Model {
	return newModel(schema, name, TpTable, version)
}

/**
* NewModel
* @param schema *Schema, name string, version int
* @return *Model
**/
func NewModel(schema *Schema, name string, version int) *Model {
	return newModel(schema, name, TpModel, version)
}

/**
* NewCollection
* @param schema *Schema, name string, version int
* @return *Model
**/
func NewCollection(schema *Schema, name string, version int) *Model {
	result := newModel(schema, name, TpCollection, version)
	result.DefineModel()

	return result
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

	if s.Show {
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
* Drop
**/
func (s *Model) Drop() {
	if s.Db == nil {
		return
	}

	for _, detail := range s.Details {
		detail.Drop()
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

	defaultDictoinary := func(mapa map[string][]*Dictionary, key string, value interface{}) map[string][]*Dictionary {
		dictionaries := mapa[key]
		if dictionaries == nil {
			return nil
		}

		idx := slices.IndexFunc(dictionaries, func(e *Dictionary) bool { return e.Value == value })
		if idx != -1 {
			dictionary := dictionaries[idx]
			defaultColValue(dictionary.Columns)
			if len(dictionary.Dictionaries) != 0 {
				return dictionary.Dictionaries
			}
		}

		return mapa
	}

	defaultColValue(s.Columns)

	dictionaries := s.Dictionaries
	for key, value := range data {
		dictionaries = defaultDictoinary(dictionaries, key, value)
		if dictionaries == nil {
			dictionaries = s.Dictionaries
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

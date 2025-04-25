package jdb

import (
	"encoding/json"
	"slices"
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/timezone"
)

func tableName(schema *Schema, name string) string {
	table := strs.Lowcase(name)
	if schema != nil {
		return strs.Format(`%s.%s`, strs.Lowcase(schema.Name), table)
	}

	return table
}

type Model struct {
	Db              *DB                  `json:"-"`
	Schema          *Schema              `json:"-"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdateAt        time.Time            `json:"updated_at"`
	Id              string               `json:"id"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	UseCore         bool                 `json:"use_core"`
	Table           string               `json:"table"`
	Integrity       bool                 `json:"integrity"`
	Definitions     []et.Json            `json:"definitions"`
	Columns         []*Column            `json:"-"`
	GeneratedFields []*Column            `json:"-"`
	PrimaryKeys     map[string]*Column   `json:"-"`
	ForeignKeys     map[string]*Column   `json:"-"`
	Indices         map[string]*Index    `json:"-"`
	Uniques         map[string]*Index    `json:"-"`
	RelationsTo     map[string]*Relation `json:"-"`
	Details         map[string]*Relation `json:"-"`
	Rollups         map[string]*Rollup   `json:"-"`
	History         *Relation            `json:"-"`
	Required        map[string]bool      `json:"-"`
	SystemKeyField  *Column              `json:"-"`
	StatusField     *Column              `json:"-"`
	IndexField      *Column              `json:"-"`
	SourceField     *Column              `json:"-"`
	FullTextField   *Column              `json:"-"`
	ProjectField    *Column              `json:"-"`
	Version         int                  `json:"version"`
	eventError      []EventError         `json:"-"`
	eventsInsert    []Event              `json:"-"`
	eventsUpdate    []Event              `json:"-"`
	eventsDelete    []Event              `json:"-"`
	isDebug         bool                 `json:"-"`
	isInit          bool                 `json:"-"`
}

/**
* NewModel
* @param schema *Schema, name string, version int
* @return *Model
**/
func NewModel(schema *Schema, name string, version int) *Model {
	name = Name(name)
	idx := slices.IndexFunc(schema.Db.models, func(model *Model) bool { return model.Name == name })
	if idx != -1 {
		return schema.Db.models[idx]
	}

	newModel := func() *Model {
		table := tableName(schema, name)
		now := timezone.NowTime()
		result := &Model{
			Db:              schema.Db,
			Schema:          schema,
			CreatedAt:       now,
			UpdateAt:        now,
			Id:              reg.Id("model"),
			Name:            name,
			UseCore:         schema.UseCore,
			Table:           table,
			Definitions:     make([]et.Json, 0),
			Columns:         make([]*Column, 0),
			GeneratedFields: make([]*Column, 0),
			PrimaryKeys:     make(map[string]*Column),
			ForeignKeys:     make(map[string]*Column),
			Indices:         make(map[string]*Index),
			Uniques:         make(map[string]*Index),
			RelationsTo:     make(map[string]*Relation),
			Details:         make(map[string]*Relation),
			Rollups:         make(map[string]*Rollup),
			Required:        make(map[string]bool),
			eventError:      make([]EventError, 0),
			eventsInsert:    make([]Event, 0),
			eventsUpdate:    make([]Event, 0),
			eventsDelete:    make([]Event, 0),
			Version:         version,
		}
		result.DefineEventError(EventErrorDefault)
		result.DefineEvent(EventInsert, EventInsertDefault)
		result.DefineEvent(EventUpdate, EventUpdateDefault)
		result.DefineEvent(EventDelete, EventDeleteDefault)

		schema.models = append(schema.models, result)
		schema.Db.models = append(schema.Db.models, result)
		return result
	}

	if !schema.UseCore || !schema.Db.isInit {
		return newModel()
	}

	var result *Model
	err := schema.Db.Load("model", name, &result)
	if err != nil {
		return nil
	}

	if result != nil {
		result.Db = schema.Db
		result.Schema = schema
		result.Columns = make([]*Column, 0)
		result.GeneratedFields = make([]*Column, 0)
		result.PrimaryKeys = make(map[string]*Column)
		result.ForeignKeys = make(map[string]*Column)
		result.Indices = make(map[string]*Index)
		result.Uniques = make(map[string]*Index)
		result.RelationsTo = make(map[string]*Relation)
		result.Details = make(map[string]*Relation)
		result.Rollups = make(map[string]*Rollup)
		result.History = nil
		result.Required = make(map[string]bool)
		// event
		result.eventError = make([]EventError, 0)
		result.eventsInsert = make([]Event, 0)
		result.eventsUpdate = make([]Event, 0)
		result.eventsDelete = make([]Event, 0)
		result.DefineEventError(EventErrorDefault)
		result.DefineEvent(EventInsert, EventInsertDefault)
		result.DefineEvent(EventUpdate, EventUpdateDefault)
		result.DefineEvent(EventDelete, EventDeleteDefault)
		// define columns
		for _, definition := range result.Definitions {
			args := definition.Array("args")
			tp := definition.Int("tp")
			result.defineColumns(tp, args...)
		}

		schema.models = append(schema.models, result)
		schema.Db.models = append(schema.Db.models, result)
		return result
	}

	return newModel()
}

/**
* Describe
* @return et.Json
**/
func (s *Model) Describe() et.Json {
	definition, err := json.Marshal(s)
	if err != nil {
		return et.Json{}
	}

	result := et.Json{}
	err = json.Unmarshal(definition, &result)
	if err != nil {
		return et.Json{}
	}

	result.Set("columns", s.Columns)
	result.Set("generated_fields", s.GeneratedFields)
	result.Set("primary_keys", s.PrimaryKeys)
	result.Set("foreign_keys", s.ForeignKeys)
	result.Set("indices", s.Indices)
	result.Set("uniques", s.Uniques)
	result.Set("relations_to", s.RelationsTo)
	result.Set("details", s.Details)
	result.Set("rollups", s.Rollups)
	result.Set("history", s.History)
	result.Set("required", s.Required)
	result.Set("system_key_field", s.SystemKeyField)
	result.Set("status_field", s.StatusField)
	result.Set("index_field", s.IndexField)
	result.Set("source_field", s.SourceField)
	result.Set("full_text_field", s.FullTextField)
	result.Set("project_field", s.ProjectField)

	return result
}

/**
* Save
* @return error
**/
func (s *Model) Save() error {
	if !s.UseCore || !s.Db.isInit {
		return nil
	}

	definition, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = s.Db.upsertModel("model", s.Name, s.Version, definition)
	if err != nil {
		return err
	}

	s.isInit = true

	return nil
}

/**
* Init
* @return error
**/
func (s *Model) Init() error {
	if !s.UseCore || s.isInit {
		return nil
	}

	if s.SourceField != nil {
		idx := s.SourceField.Idx()
		if idx != len(s.Columns)-1 && idx > -1 {
			s.Columns = append(s.Columns[:idx], s.Columns[idx+1:]...)
			s.Columns = append(s.Columns, s.SourceField)
		}
	}

	if s.IndexField != nil {
		idx := s.IndexField.Idx()
		if idx != len(s.Columns)-1 && idx > -1 {
			s.Columns = append(s.Columns[:idx], s.Columns[idx+1:]...)
			s.Columns = append(s.Columns, s.IndexField)
		}
	}

	if s.SystemKeyField != nil {
		idx := s.SystemKeyField.Idx()
		if idx != len(s.Columns)-1 && idx > -1 {
			s.Columns = append(s.Columns[:idx], s.Columns[idx+1:]...)
			s.Columns = append(s.Columns, s.SystemKeyField)
		}
	}

	err := s.Db.LoadModel(s)
	if err != nil {
		return err
	}

	err = s.Save()
	if err != nil {
		return err
	}

	return nil
}

/**
* Drop
**/
func (s *Model) Drop() {
	if s.Db == nil {
		return
	}

	for _, detail := range s.Details {
		model := detail.With
		if model != nil {
			model.Drop()
		}
	}

	s.Db.DropModel(s)
}

/**
* IsDebug
* @return bool
**/
func (s *Model) IsDebug() bool {
	return s.isDebug
}

/**
* GetFrom
* @return *QlFrom
**/
func (s *Model) GetFrom() *QlFrom {
	return &QlFrom{Model: s}
}

/**
* GetId
* @param id string
* @return string
**/
func (s *Model) GetId(id string) string {
	return reg.GetId(s.Table, id)
}

/**
* GenId
* @return string
**/
func (s *Model) GenId() string {
	return reg.Id(s.Table)
}

/**
* GetSerie
* @return int64, error
**/
func (s *Model) GetSerie() (int64, error) {
	return s.Db.GetSerie(s.Table)
}

/**
* GetCode
* @param tag, prefix string
* @return string, error
**/
func (s *Model) GetCode(tag, prefix string) (string, error) {
	return s.Db.GetCode(tag, prefix)
}

/**
* SetSerie
* @param tag string, val int64
* @return int64, error
**/
func (s *Model) SetSerie(tag string, val int64) (int64, error) {
	return s.Db.SetSerie(tag, val)
}

/**
* sourceIdx
* @return int
**/
func (s *Model) sourceIdx() int {
	if s.SourceField == nil {
		return -1
	}

	return s.SourceField.Idx()
}

/**
* Debug
* @return *Model
**/
func (s *Model) Debug() *Model {
	s.isDebug = true

	return s
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
* GetColumnsArray
* @param names ...string
* @return []string
**/
func (s *Model) GetColumnsArray(names ...string) []string {
	result := []string{}
	for _, name := range names {
		if col := s.GetColumn(name); col != nil {
			result = append(result, col.Name)
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
	if col != nil {
		return col.GetField()
	}

	if s.Integrity {
		return nil
	}

	if s.SourceField == nil {
		return nil
	}

	result := s.defineAtribute(name, TypeDataText)

	return result.GetField()
}

/**
* GetField
* @param name string
* @return *Field
**/
func (s *Model) GetField(name string) *Field {
	list := strs.Split(name, ":")
	alias := ""
	if len(list) > 1 {
		name = list[0]
		alias = list[1]
	}

	list = strs.Split(name, ".")
	switch len(list) {
	case 1:
		result := s.getField(list[0])
		if alias != "" {
			result.Alias = alias
		}

		return result
	case 2:
		if !strs.Same(s.Name, list[0]) {
			return nil
		}

		result := s.getField(list[1])
		if alias != "" {
			result.Alias = alias
		}

		return result
	case 3:
		table := strs.Format(`%s.%s`, list[0], list[1])
		if !strs.Same(s.Table, table) {
			return nil
		}

		result := s.getField(list[2])
		if alias != "" {
			result.Alias = alias
		}

		return result
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
	for _, col := range s.PrimaryKeys {
		result = append(result, col)
	}

	return result
}

/**
* Where
* @param val string
* @return *Ql
**/
func (s *Model) Where(val string) *Ql {
	result := From(s)
	if s.SourceField != nil {
		result.TypeSelect = Data
	}

	return result.Where(val)
}

/**
* Query
* @param params et.Json
* @return interface{}, error
**/
func (s *Model) Query(params et.Json) (interface{}, error) {
	return From(s).
		Query(params)
}

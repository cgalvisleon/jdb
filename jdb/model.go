package jdb

import (
	"encoding/json"
	"slices"
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/timezone"
)

type Model struct {
	Db                 *DB                      `json:"-"`
	Schema             *Schema                  `json:"-"`
	CreatedAt          time.Time                `json:"created_at"`
	UpdateAt           time.Time                `json:"updated_at"`
	Id                 string                   `json:"id"`
	Name               string                   `json:"name"`
	Description        string                   `json:"description"`
	UseCore            bool                     `json:"use_core"`
	Integrity          bool                     `json:"integrity"`
	Definitions        []et.Json                `json:"definitions"`
	Columns            []*Column                `json:"-"`
	PrimaryKeys        map[string]*Column       `json:"-"`
	ForeignKeys        map[string]*Relation     `json:"-"`
	Indices            map[string]*Index        `json:"-"`
	Uniques            map[string]*Index        `json:"-"`
	RelationsTo        map[string]*Relation     `json:"-"`
	Details            map[string]*Relation     `json:"-"`
	Rollups            map[string]*Rollup       `json:"-"`
	History            *Relation                `json:"-"`
	Required           map[string]bool          `json:"-"`
	CreatedAtField     *Column                  `json:"-"`
	UpdatedAtField     *Column                  `json:"-"`
	SystemKeyField     *Column                  `json:"-"`
	StatusField        *Column                  `json:"-"`
	IndexField         *Column                  `json:"-"`
	SourceField        *Column                  `json:"-"`
	FullTextField      *Column                  `json:"-"`
	ProjectField       *Column                  `json:"-"`
	Version            int                      `json:"version"`
	eventError         []EventError             `json:"-"`
	eventsInsert       []Event                  `json:"-"`
	eventsUpdate       []Event                  `json:"-"`
	eventsDelete       []Event                  `json:"-"`
	eventEmiterChannel chan event.Message       `json:"-"`
	eventsEmiter       map[string]event.Handler `json:"-"`
	IsDebug            bool                     `json:"-"`
	isLocked           bool                     `json:"-"`
	isInit             bool                     `json:"-"`
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
		now := timezone.NowTime()
		result := &Model{
			Db:                 schema.Db,
			Schema:             schema,
			CreatedAt:          now,
			UpdateAt:           now,
			Id:                 reg.GenId("model"),
			Name:               name,
			UseCore:            schema.UseCore,
			Definitions:        make([]et.Json, 0),
			Columns:            make([]*Column, 0),
			PrimaryKeys:        make(map[string]*Column),
			ForeignKeys:        make(map[string]*Relation),
			Indices:            make(map[string]*Index),
			Uniques:            make(map[string]*Index),
			RelationsTo:        make(map[string]*Relation),
			Details:            make(map[string]*Relation),
			Rollups:            make(map[string]*Rollup),
			Required:           make(map[string]bool),
			eventEmiterChannel: make(chan event.Message),
			eventsEmiter:       make(map[string]event.Handler),
			eventError:         make([]EventError, 0),
			eventsInsert:       make([]Event, 0),
			eventsUpdate:       make([]Event, 0),
			eventsDelete:       make([]Event, 0),
			Version:            version,
		}
		result.DefineEventError(eventErrorDefault)
		result.DefineEvent(EventInsert, eventInsertDefault)
		result.DefineEvent(EventUpdate, eventUpdateDefault)
		result.DefineEvent(EventDelete, eventDeleteDefault)

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
		result.PrimaryKeys = make(map[string]*Column)
		result.ForeignKeys = make(map[string]*Relation)
		result.Indices = make(map[string]*Index)
		result.Uniques = make(map[string]*Index)
		result.RelationsTo = make(map[string]*Relation)
		result.Details = make(map[string]*Relation)
		result.Rollups = make(map[string]*Rollup)
		result.History = nil
		result.Required = make(map[string]bool)
		/* Event */
		result.eventEmiterChannel = make(chan event.Message)
		result.eventsEmiter = make(map[string]event.Handler)
		result.eventError = make([]EventError, 0)
		result.eventsInsert = make([]Event, 0)
		result.eventsUpdate = make([]Event, 0)
		result.eventsDelete = make([]Event, 0)
		result.DefineEventError(eventErrorDefault)
		result.DefineEvent(EventInsert, eventInsertDefault)
		result.DefineEvent(EventUpdate, eventUpdateDefault)
		result.DefineEvent(EventDelete, eventDeleteDefault)

		/* Define columns */
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
* Serialize
* @return []byte, error
**/
func (s *Model) serialize() ([]byte, error) {
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
func (s *Model) Describe() et.Json {
	definition, err := s.serialize()
	if err != nil {
		return et.Json{}
	}

	result := et.Json{}
	err = json.Unmarshal(definition, &result)
	if err != nil {
		return et.Json{}
	}

	columns := make([]et.Json, 0)
	for _, column := range s.Columns {
		columns = append(columns, column.Describe())
	}

	delete(result, "definitions")
	result["kind"] = "model"
	result["columns"] = columns
	result["primary_keys"] = s.PrimaryKeys
	result["foreign_keys"] = s.ForeignKeys
	result["indices"] = s.Indices
	result["uniques"] = s.Uniques
	result["relations_to"] = s.RelationsTo
	result["details"] = s.Details
	result["rollups"] = s.Rollups
	result["history"] = s.History
	result["required"] = s.Required
	result["system_key_field"] = s.SystemKeyField
	result["status_field"] = s.StatusField
	result["index_field"] = s.IndexField
	result["source_field"] = s.SourceField
	result["full_text_field"] = s.FullTextField
	result["project_field"] = s.ProjectField

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

	definition, err := s.serialize()
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
	go func() {
		for message := range s.eventEmiterChannel {
			s.eventEmiter(message)
		}
	}()

	if !s.UseCore || s.isInit {
		return nil
	}

	if s.SourceField != nil {
		idx := s.SourceField.idx()
		if idx != len(s.Columns)-1 && idx != -1 {
			s.Columns = append(s.Columns[:idx], s.Columns[idx+1:]...)
			s.Columns = append(s.Columns, s.SourceField)
		}
	}

	if s.SystemKeyField != nil {
		idx := s.SystemKeyField.idx()
		if idx != len(s.Columns)-1 && idx > -1 {
			s.Columns = append(s.Columns[:idx], s.Columns[idx+1:]...)
			s.Columns = append(s.Columns, s.SystemKeyField)
		}
	}

	if s.IndexField != nil {
		idx := s.IndexField.idx()
		if idx != len(s.Columns)-1 && idx != -1 {
			s.Columns = append(s.Columns[:idx], s.Columns[idx+1:]...)
			s.Columns = append(s.Columns, s.IndexField)
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
	return reg.GetId(s.Name, id)
}

/**
* GenId
* @return string
**/
func (s *Model) GenId() string {
	return reg.GenId(s.Name)
}

/**
* GetSerie
* @return int64, error
**/
func (s *Model) GetSerie() (int64, error) {
	return s.Db.GetSerie(s.Name)
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

	return s.SourceField.idx()
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
* getColumn
* @param name string
* @return *Column
**/
func (s *Model) getColumn(name string) *Column {
	for _, col := range s.Columns {
		if col.Name == name {
			return col
		}
	}

	return nil
}

/**
* getColumns
* @param name string
* @return *Column
*
 */
func (s *Model) getColumns(names ...string) []*Column {
	result := []*Column{}
	for _, name := range names {
		if col := s.getColumn(name); col != nil {
			result = append(result, col)
		}
	}

	return result
}

/**
* getColumnsArray
* @param names ...string
* @return []string
**/
func (s *Model) getColumnsArray(names ...string) []string {
	result := []string{}
	for _, name := range names {
		if col := s.getColumn(name); col != nil {
			result = append(result, col.Name)
		}
	}

	return result
}

/**
* getField
* @param name string, isCreate bool
* @return *Field
**/
func (s *Model) getField(name string, isCreate bool) *Field {
	getField := func(name string, isCreate bool) *Field {
		col := s.getColumn(name)
		if col != nil {
			return col.GetField()
		}

		if s.Integrity {
			return nil
		}

		if s.SourceField == nil {
			return nil
		}

		if !isCreate {
			return nil
		}

		result := newAtribute(s, name, TypeDataText)

		return result.GetField()
	}

	list := strs.Split(name, ":")
	alias := ""
	if len(list) > 1 {
		name = list[0]
		alias = list[1]
	}

	list = strs.Split(name, ".")
	switch len(list) {
	case 1:
		result := getField(list[0], isCreate)
		if result != nil && alias != "" {
			result.Alias = alias
		}

		return result
	case 2:
		if !strs.Same(s.Name, list[0]) {
			return nil
		}

		result := getField(list[1], isCreate)
		if result != nil && alias != "" {
			result.Alias = alias
		}

		return result
	case 3:
		if !strs.Same(s.Schema.Name, list[0]) {
			return nil
		}

		if !strs.Same(s.Name, list[1]) {
			return nil
		}

		result := getField(list[2], isCreate)
		if result != nil && alias != "" {
			result.Alias = alias
		}

		return result
	default:
		return nil
	}
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
* QueryTx
* @param tx *Tx, params et.Json
* @return interface{}, error
**/
func (s *Model) QueryTx(tx *Tx, params et.Json) (interface{}, error) {
	return From(s).
		queryTx(tx, params)
}

/**
* Query
* @param params et.Json
* @return interface{}, error
**/
func (s *Model) Query(params et.Json) (interface{}, error) {
	return From(s).
		queryTx(nil, params)
}

package jdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/timezone"
)

var (
	ErrNotInserted = errors.New("not inserted")
	ErrNotUpdated  = errors.New("not updated")
)

type Model struct {
	Db                 *DB                      `json:"-"`
	schema             *Schema                  `json:"-"`
	Schema             string                   `json:"schema"`
	CreatedAt          time.Time                `json:"created_at"`
	UpdateAt           time.Time                `json:"updated_at"`
	Id                 string                   `json:"id"`
	Name               string                   `json:"name"`
	Description        string                   `json:"description"`
	UseCore            bool                     `json:"use_core"`
	Integrity          bool                     `json:"integrity"`
	Definitions        et.Json                  `json:"definitions"`
	Columns            []*Column                `json:"-"`
	PrimaryKeys        map[string]*Column       `json:"-"`
	ForeignKeys        map[string]*Relation     `json:"-"`
	Indices            map[string]*Index        `json:"-"`
	Uniques            map[string]*Index        `json:"-"`
	RelationsTo        map[string]*Relation     `json:"-"`
	Details            map[string]*Relation     `json:"-"`
	Joins              map[string]*Join         `json:"-"`
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
	isCore             bool                     `json:"-"`
}

/**
* NewModel
* @param schema *Schema, name string, version int
* @return *Model
**/
func NewModel(schema *Schema, name string, version int) *Model {
	name = Name(name)
	idx := slices.IndexFunc(schema.Db.models, func(e *Model) bool { return e.Name == name })
	if idx != -1 {
		return schema.Db.models[idx]
	}

	newModel := func() *Model {
		if !schema.isCore {
			console.Logf("model", `Model %s new`, name)
		}

		now := timezone.NowTime()
		result := &Model{
			Db:                 schema.Db,
			schema:             schema,
			Schema:             schema.Name,
			CreatedAt:          now,
			UpdateAt:           now,
			Id:                 reg.GenId("model"),
			Name:               name,
			UseCore:            schema.UseCore,
			Definitions:        et.Json{},
			Columns:            make([]*Column, 0),
			PrimaryKeys:        make(map[string]*Column),
			ForeignKeys:        make(map[string]*Relation),
			Indices:            make(map[string]*Index),
			Uniques:            make(map[string]*Index),
			RelationsTo:        make(map[string]*Relation),
			Details:            make(map[string]*Relation),
			Joins:              make(map[string]*Join),
			Required:           make(map[string]bool),
			eventEmiterChannel: make(chan event.Message),
			eventsEmiter:       make(map[string]event.Handler),
			eventError:         make([]EventError, 0),
			eventsInsert:       make([]Event, 0),
			eventsUpdate:       make([]Event, 0),
			eventsDelete:       make([]Event, 0),
			Version:            version,
			isCore:             schema.isCore,
		}
		result.DefineEventError(eventErrorDefault)
		result.DefineEvent(EventInsert, eventInsertDefault)
		result.DefineEvent(EventUpdate, eventUpdateDefault)
		result.DefineEvent(EventDelete, eventDeleteDefault)

		schema.addModel(result)
		return result
	}

	if !schema.UseCore || !schema.Db.isInit {
		return newModel()
	}

	var result *Model
	err := schema.Db.Load("model", name, &result)
	if err != nil {
		return newModel()
	}

	result, err = loadModel(schema, result)
	if err != nil {
		result = newModel()
	}

	return result
}

/**
* loadModel
* @param schema *Schema, model *Model
* @return *Model, error
**/
func loadModel(schema *Schema, model *Model) (*Model, error) {
	idx := slices.IndexFunc(schema.Db.models, func(e *Model) bool { return e.Name == model.Name })
	if idx != -1 {
		return schema.Db.models[idx], nil
	}

	if !schema.isCore {
		console.Logf("model", `Model %s load`, model.Name)
	}

	schema.addModel(model)
	model.schema = schema
	model.Db = schema.Db
	model.Schema = schema.Name
	model.Columns = make([]*Column, 0)
	model.PrimaryKeys = make(map[string]*Column)
	model.ForeignKeys = make(map[string]*Relation)
	model.Indices = make(map[string]*Index)
	model.Uniques = make(map[string]*Index)
	model.RelationsTo = make(map[string]*Relation)
	model.Details = make(map[string]*Relation)
	model.Required = make(map[string]bool)
	/* Event */
	model.eventEmiterChannel = make(chan event.Message)
	model.eventsEmiter = make(map[string]event.Handler)
	model.eventError = make([]EventError, 0)
	model.eventsInsert = make([]Event, 0)
	model.eventsUpdate = make([]Event, 0)
	model.eventsDelete = make([]Event, 0)
	model.DefineEventError(eventErrorDefault)
	model.DefineEvent(EventInsert, eventInsertDefault)
	model.DefineEvent(EventUpdate, eventUpdateDefault)
	model.DefineEvent(EventDelete, eventDeleteDefault)
	/* Define columns */
	for name := range model.Definitions {
		definition := model.Definitions.Json(name)
		args := definition.Array("args")
		tp := definition.Int("tp")
		model.defineColumns(tp, args...)
	}

	return model, nil
}

/**
* LoadModel
* @param db *DB, name string
* @return *Model, error
**/
func LoadModel(db *DB, name string) (*Model, error) {
	name = Name(name)
	idx := slices.IndexFunc(db.models, func(e *Model) bool { return e.Name == name })
	if idx != -1 {
		return db.models[idx], nil
	}

	var result *Model
	err := db.Load("model", name, &result)
	if err != nil {
		return nil, err
	}

	if result != nil {
		schema, err := LoadSchema(db, result.Schema)
		if err != nil {
			return nil, err
		}

		return loadModel(schema, result)
	}

	return result, nil
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
			s.moveColumnToEnd(s.SourceField, idx)
		}
	}

	if s.SystemKeyField != nil {
		idx := s.SystemKeyField.idx()
		if idx != len(s.Columns)-1 && idx != -1 {
			s.moveColumnToEnd(s.SystemKeyField, idx)
		}
	}

	if s.IndexField != nil {
		idx := s.IndexField.idx()
		if idx != len(s.Columns)-1 && idx != -1 {
			s.moveColumnToEnd(s.IndexField, idx)
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
* GetWhereByPk
* @param data et.Json
* @return et.Json
**/
func (s *Model) GetWhereByPk(data et.Json) (et.Json, error) {
	result := et.Json{}
	for name := range s.PrimaryKeys {
		val := data.Get(name)
		if val == nil {
			return et.Json{}, mistake.Newf(MSG_FIELD_REQUIRED_RELATION, name, s.Name)
		}

		col := s.getColumn(name)
		if col != nil && col.IsKeyfield {
			vs := fmt.Sprintf(`%v`, val)
			val = s.GetId(vs)
		}

		result[name] = et.Json{
			"eq": val,
		}
	}

	return result, nil
}

/**
* GetWhereByRequired
* @param data et.Json
* @return et.Json
**/
func (s *Model) GetWhereByRequired(data et.Json) (et.Json, error) {
	result := et.Json{}
	for name := range s.Required {
		val := data.Get(name)

		if val == nil {
			return et.Json{}, mistake.Newf(MSG_FIELD_REQUIRED_RELATION, name, s.Name)
		}

		col := s.getColumn(name)
		if col != nil && col.IsKeyfield {
			vs := fmt.Sprintf(`%v`, val)
			val = s.GetId(vs)
		}

		result[name] = et.Json{
			"eq": val,
		}
	}

	return result, nil
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
* addColumn
* @param column *Column
**/
func (s *Model) addColumn(column *Column) {
	idx := slices.IndexFunc(s.Columns, func(e *Column) bool { return e.Name == column.Name })
	if idx == -1 {
		s.Columns = append(s.Columns, column)
	}
}

/**
* addColumnIdx
* @param column *Column, idx int
**/
func (s *Model) addColumnToIdx(column *Column, idx int) {
	if idx != -1 {
		s.Columns = append(s.Columns[:idx], append([]*Column{column}, s.Columns[idx:]...)...)
	}
}

/**
* moveColumnToEnd
* @param column *Column, idx int
**/
func (s *Model) moveColumnToEnd(column *Column, idx int) {
	s.Columns = append(s.Columns[:idx], s.Columns[idx+1:]...)
	s.addColumn(column)
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
		if !strs.Same(s.Schema, list[0]) {
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
		result.TypeSelect = Source
	}

	return result.Where(val)
}

/**
* Counted
* @return int, error
**/
func (s *Model) Counted() (int, error) {
	all, err := From(s).
		Counted()
	if err != nil {
		return 0, err
	}

	return all, nil
}

/**
* QueryTx
* @param tx *Tx, params et.Json
* @return et.Json, error
**/
func (s *Model) QueryTx(tx *Tx, params et.Json) (et.Json, error) {
	return From(s).
		queryTx(tx, params)
}

/**
* Query
* @param params et.Json
* @return et.Json, error
**/
func (s *Model) Query(params et.Json) (et.Json, error) {
	return s.QueryTx(nil, params)
}

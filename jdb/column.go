package jdb

import (
	"encoding/gob"
	"encoding/json"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/timezone"
)

type TypeColumn int

const (
	TpColumn TypeColumn = iota
	TpAtribute
	TpCalc
	TpRelatedTo
	TpRollup
	TpConcurrent
)

func (s TypeColumn) Str() string {
	switch s {
	case TpColumn:
		return "column"
	case TpAtribute:
		return "attribute"
	case TpCalc:
		return "calc"
	case TpRelatedTo:
		return "related_to"
	case TpRollup:
		return "rollup"
	case TpConcurrent:
		return "concurrent"
	}

	return "attribute"
}

/**
* StrsToTypeColumn
* @param strs string
* @return TypeColumn
**/
func StrsToTypeColumn(strs string) TypeColumn {
	switch strs {
	case "column":
		return TpColumn
	case "attribute":
		return TpAtribute
	case "calc":
		return TpCalc
	case "related_to":
		return TpRelatedTo
	case "rollup":
		return TpRollup
	case "concurrent":
		return TpConcurrent
	}

	return TpAtribute
}

type TypeData int

const (
	TypeDataText TypeData = iota
	TypeDataMemo
	TypeDataShortText
	TypeDataKey
	TypeDataNumber
	TypeDataInt
	TypeDataPrecision
	TypeDataDateTime
	TypeDataCheckbox
	TypeDataBytes
	/* Jsonb */
	TypeDataObject
	TypeDataSelect
	TypeDataMultiSelect
	TypeDataGeometry
	/* Fulltext */
	TypeDataFullText
	/* Objects */
	TypeDataState
	TypeDataUser
	TypeDataFilesMedia
	TypeDataUrl
	TypeDataEmail
	TypeDataPhone
	TypeDataAddress
	TypeDataRelation
	TypeDataRollup
	/* None */
	TypeDataNone
)

func (s TypeData) DefaultValue() interface{} {
	switch s {
	case TypeDataNumber:
		return 0.0
	case TypeDataInt:
		return 0
	case TypeDataCheckbox:
		return false
	case TypeDataSelect:
		return ""
	case TypeDataMultiSelect:
		return []string{}
	case TypeDataMemo:
		return ""
	case TypeDataObject:
		return et.Json{}
	case TypeDataGeometry:
		return et.Json{
			"type":        "Point",
			"coordinates": []float64{0.00, 0.00},
		}
	}

	return ""
}

func (s TypeData) Str() string {
	switch s {
	case TypeDataText:
		return "text"
	case TypeDataMemo:
		return "memo"
	case TypeDataShortText:
		return "short-text"
	case TypeDataKey:
		return "key"
	case TypeDataNumber:
		return "number"
	case TypeDataInt:
		return "int"
	case TypeDataPrecision:
		return "precision"
	case TypeDataDateTime:
		return "date-time"
	case TypeDataBytes:
		return "bytes"
	case TypeDataObject:
		return "object"
	case TypeDataFullText:
		return "full-text"
	case TypeDataSelect:
		return "select"
	case TypeDataMultiSelect:
		return "multi-select"
	case TypeDataGeometry:
		return "location"
	case TypeDataCheckbox:
		return "checkbox"
	case TypeDataState:
		return "status"
	case TypeDataUser:
		return "person"
	case TypeDataFilesMedia:
		return "files-media"
	case TypeDataUrl:
		return "url"
	case TypeDataEmail:
		return "email"
	case TypeDataPhone:
		return "phone"
	case TypeDataAddress:
		return "address"
	case TypeDataRelation:
		return "relation"
	case TypeDataRollup:
		return "rollup"
	default:
		return "none"
	}
}

func StrsToTypeData(strs string) TypeData {
	switch strs {
	case "text":
		return TypeDataText
	case "memo":
		return TypeDataMemo
	case "number":
		return TypeDataNumber
	case "select":
		return TypeDataSelect
	case "multi-select":
		return TypeDataMultiSelect
	case "date-time":
		return TypeDataDateTime
	case "checkbox":
		return TypeDataCheckbox
	case "status":
		return TypeDataState
	case "person":
		return TypeDataUser
	case "files-media":
		return TypeDataFilesMedia
	case "url":
		return TypeDataUrl
	case "email":
		return TypeDataEmail
	case "phone":
		return TypeDataPhone
	case "address":
		return TypeDataAddress
	case "location":
		return TypeDataGeometry
	case "relation":
		return TypeDataRelation
	case "rollup":
		return TypeDataRollup
	}

	return TypeDataText
}

func StrToKindType(strs string) (TypeColumn, TypeData) {
	switch strs {
	case "text":
		return TpAtribute, TypeDataText
	case "memo":
		return TpAtribute, TypeDataMemo
	case "number":
		return TpAtribute, TypeDataNumber
	case "select":
		return TpAtribute, TypeDataSelect
	case "multi-select":
		return TpAtribute, TypeDataMultiSelect
	case "date-time":
		return TpAtribute, TypeDataDateTime
	case "checkbox":
		return TpAtribute, TypeDataCheckbox
	case "status":
		return TpColumn, TypeDataState
	case "person":
		return TpColumn, TypeDataUser
	case "files-media":
		return TpColumn, TypeDataFilesMedia
	case "url":
		return TpColumn, TypeDataUrl
	case "email":
		return TpColumn, TypeDataEmail
	case "phone":
		return TpColumn, TypeDataPhone
	case "address":
		return TpAtribute, TypeDataAddress
	case "location":
		return TpColumn, TypeDataGeometry
	case "relation":
		return TpRelatedTo, TypeDataRelation
	case "rollup":
		return TpRollup, TypeDataRollup
	}

	return TpAtribute, TypeDataText
}

type ColumnFields struct {
	Key        string
	Index      string
	Source     string
	ProjectId  string
	CreatedAt  string
	UpdatedAt  string
	StatusId   string
	SystemId   string
	CreatedTo  string
	UpdatedTo  string
	Fulltext   string
	Historical string
	Checked    string
}

/**
* Json
* @return et.Json
**/
func (s *ColumnFields) Json() et.Json {
	return et.Json{
		"key":        s.Key,
		"index":      s.Index,
		"source":     s.Source,
		"project_id": s.ProjectId,
		"created_at": s.CreatedAt,
		"updated_at": s.UpdatedAt,
		"status_id":  s.StatusId,
		"system_id":  s.SystemId,
		"created_to": s.CreatedTo,
		"updated_to": s.UpdatedTo,
		"fulltext":   s.Fulltext,
		"historical": s.Historical,
		"checked":    s.Checked,
	}
}

var (
	cf         *ColumnFields
	KEY        string
	PRIMARYKEY string
	INDEX      string
	SOURCE     string
	PROJECT_ID string
	CREATED_AT string
	UPDATED_AT string
	STATUS_ID  string
	SYSID      string
	CREATED_TO string
	UPDATED_TO string
	FULLTEXT   string
	HISTORYCAL string
	CHECKED    string
)

func init() {
	cf = &ColumnFields{
		Key:        "id",
		Index:      "idx",
		Source:     "source",
		ProjectId:  "project_id",
		CreatedAt:  "created_at",
		UpdatedAt:  "updated_at",
		StatusId:   "status_id",
		SystemId:   "jdbid",
		CreatedTo:  "created_to",
		UpdatedTo:  "updated_to",
		Fulltext:   "fulltext",
		Historical: "historical",
		Checked:    "checked",
	}

	SetColumnFields(cf)
}

/**
* SetColumnFields
* @param fields *ColumnFields
**/
func SetColumnFields(fields *ColumnFields) {
	cf = fields
	KEY = cf.Key
	PRIMARYKEY = cf.Key
	INDEX = cf.Index
	SOURCE = cf.Source
	PROJECT_ID = cf.ProjectId
	CREATED_AT = cf.CreatedAt
	UPDATED_AT = cf.UpdatedAt
	STATUS_ID = cf.StatusId
	SYSID = cf.SystemId
	CREATED_TO = cf.CreatedTo
	UPDATED_TO = cf.UpdatedTo
	FULLTEXT = cf.Fulltext
	HISTORYCAL = cf.Historical
	CHECKED = cf.Checked
}

type Relation struct {
	With            *Model            `json:"-"`
	Fk              map[string]string `json:"fk"`
	Limit           int               `json:"rows"`
	OnDeleteCascade bool              `json:"on_delete_cascade"`
	OnUpdateCascade bool              `json:"on_update_cascade"`
	IsMultiSelect   bool              `json:"is_multi_select"`
}

/**
* GetWhere
* @param from et.Json
* @return et.Json
**/
func (s *Relation) GetWhere(from et.Json) et.Json {
	result := et.Json{}
	for fkn, pkn := range s.Fk {
		fk := from.Get(fkn)
		if fk == nil {
			continue
		}

		result[pkn] = et.Json{
			"eq": fk,
		}
	}

	return result
}

/**
* Serialize
* @return []byte, error
**/
func (s *Relation) Serialize() ([]byte, error) {
	result, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

/**
* describe
* @return et.Json
**/
func (s *Relation) describe() et.Json {
	definition, err := s.Serialize()
	if err != nil {
		return et.Json{}
	}

	result := et.Json{}
	err = json.Unmarshal(definition, &result)
	if err != nil {
		return et.Json{}
	}

	result["with"] = s.With.Name

	return result
}

type ShowRollup int

const (
	ShowAtrib ShowRollup = iota
	ShowObject
)

/**
* Str
* @return string
**/
func (s ShowRollup) Str() string {
	switch s {
	case ShowObject:
		return "objects"
	default:
		return "attribs"
	}
}

type Rollup struct {
	With   *Model            `json:"-"`
	Fk     map[string]string `json:"fk"`
	Fields []string          `json:"fields"`
	Show   ShowRollup        `json:"show"`
}

/**
* Where
* @param from et.Json
* @return et.Json
**/
func (s *Rollup) Where(from et.Json) et.Json {
	result := et.Json{}
	for fkn, pkn := range s.Fk {
		fk := from.Str(fkn)
		result[pkn] = et.Json{
			"eq": fk,
		}
	}

	return result
}

/**
* Serialize
* @return []byte, error
**/
func (s *Rollup) Serialize() ([]byte, error) {
	result, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

/**
* describe
* @return et.Json
**/
func (s *Rollup) describe() et.Json {
	definition, err := s.Serialize()
	if err != nil {
		return et.Json{}
	}

	result := et.Json{}
	err = json.Unmarshal(definition, &result)
	if err != nil {
		return et.Json{}
	}

	result["with"] = s.With.Name
	result["show"] = s.Show.Str()
	return result
}

type Join struct {
	On   et.Json  `json:"on"`
	Type TypeJoin `json:"type"`
}

type FullText struct {
	Language string   `json:"language"`
	Columns  []string `json:"columns"`
}

/**
* Describe
* @return et.Json
**/
func (s *FullText) describe() et.Json {
	result := et.Json{
		"language": s.Language,
		"columns":  s.Columns,
	}

	return result
}

type Column struct {
	Model        *Model       `json:"-"`
	Source       *Column      `json:"-"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	TypeColumn   TypeColumn   `json:"type_column"`
	TypeData     TypeData     `json:"type_data"`
	Default      interface{}  `json:"default"`
	IsKeyfield   bool         `json:"is_keyfield"`
	Max          float64      `json:"max"`
	Min          float64      `json:"min"`
	Hidden       bool         `json:"hidden"`
	Detail       *Relation    `json:"detail"`
	Rollup       *Rollup      `json:"rollup"`
	FullText     *FullText    `json:"fulltext"`
	Values       interface{}  `json:"values"`
	CalcFunction DataFunction `json:"-"`
}

func init() {
	gob.Register(&Column{})
}

/**
* newColumn
* @param model *Model, name string, description string, typeColumn TypeColumn, typeData TypeData, def interface{}
* @return *Column
**/
func newColumn(model *Model, name string, description string, typeColumn TypeColumn, typeData TypeData, def interface{}) *Column {
	return &Column{
		Model:       model,
		Name:        name,
		Description: description,
		TypeColumn:  typeColumn,
		TypeData:    typeData,
		Default:     def,
		Max:         0,
		Min:         0,
		Values:      "",
	}
}

/**
* newAtribute
* @param model *Model, name string, typeData TypeData
* @return *Column
**/
func newAtribute(model *Model, name string, typeData TypeData) *Column {
	if model.SourceField == nil {
		return nil
	}

	result := model.getColumn(name)
	if result != nil {
		return result
	}

	result = newColumn(model, name, "", TpAtribute, typeData, typeData.DefaultValue())
	result.Source = model.SourceField

	return result
}

/**
* Describe
* @return et.Json
**/
func (s *Column) Describe() et.Json {
	detail := et.Json{}
	if s.Detail != nil {
		detail = s.Detail.describe()
	}

	rollup := et.Json{}
	if s.Rollup != nil {
		rollup = s.Rollup.describe()
	}

	source := ""
	if s.Source != nil {
		source = s.Source.Name
	}

	fulltext := et.Json{}
	if s.FullText != nil {
		fulltext = s.FullText.describe()
	}

	result := et.Json{
		"source":      source,
		"name":        s.Name,
		"description": s.Description,
		"type_column": s.TypeColumn.Str(),
		"type_data":   s.TypeData.Str(),
		"default":     s.Default,
		"max":         s.Max,
		"min":         s.Min,
		"hidden":      s.Hidden,
		"detail":      detail,
		"rollup":      rollup,
		"fulltext":    fulltext,
		"values":      s.Values,
	}

	return result
}

/**
* Idx
* @return int
**/
func (s *Column) idx() int {
	if s.Model == nil {
		return -1
	}

	return slices.IndexFunc(s.Model.Columns, func(e *Column) bool { return e.Name == s.Name })
}

/**
* SetValue
* @param value interface{}
**/
func (s *Column) SetValue(value interface{}) *Column {
	s.Values = value
	return s
}

/**
* SetDefaultValue
* @param value interface{}
* @return *Column
**/
func (s *Column) SetDefaultValue(value interface{}) *Column {
	s.Default = value
	return s
}

/**
* SetHidden
* @param hidden bool
* @return *Column
**/
func (s *Column) SetHidden(hidden bool) *Column {
	s.Hidden = hidden
	return s
}

/**
* SetMax
* @param max float64
* @return *Column
**/
func (s *Column) SetMax(max float64) *Column {
	s.Max = max
	return s
}

/**
* SetMin
* @param min float64
* @return *Column
**/
func (s *Column) SetMin(min float64) *Column {
	s.Min = min
	return s
}

/**
* DefaultValue
* @return interface{}
**/
func (s *Column) DefaultValue() interface{} {
	if s.TypeData == TypeDataDateTime {
		return timezone.Now()
	}

	return s.Default
}

/**
* DefaultQuote
* @return interface{}
**/
func (s *Column) DefaultQuote() interface{} {
	result := s.DefaultValue()
	result = Quote(result)

	return result
}

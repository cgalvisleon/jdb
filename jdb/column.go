package jdb

import (
	"encoding/gob"
	"encoding/json"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/timezone"
	"github.com/cgalvisleon/et/utility"
)

type TypeColumn int

const (
	TpColumn TypeColumn = iota
	TpAtribute
	TpCalc
	TpRelatedTo
	TpRollup
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
	}

	return "column"
}

type TypeData int

const (
	TypeDataText TypeData = iota
	TypeDataMemo
	TypeDataShortText
	TypeDataKey
	TypeDataState
	TypeDataInt
	TypeDataNumber
	TypeDataPrecision
	TypeDataSerie
	TypeDataIndex
	TypeDataBool
	TypeDataTime
	TypeDataBytes
	/* Special */
	TypeDataObject
	TypeDataArray
	TypeDataGeometry
	TypeDataFullText
	TypeDataSelect
	TypeDataNone
)

func (s TypeData) DefaultValue() interface{} {
	switch s {
	case TypeDataArray:
		return "[]"
	case TypeDataBool:
		return false
	case TypeDataInt:
		return 0
	case TypeDataKey:
		return "-1"
	case TypeDataState:
		return utility.ACTIVE
	case TypeDataNumber:
		return 0.0
	case TypeDataPrecision:
		return 0.0
	case TypeDataObject:
		return "{}"
	case TypeDataSerie:
		return 0
	case TypeDataIndex:
		return 0
	case TypeDataTime:
		return "NOW()"
	case TypeDataGeometry:
		return "{type: 'Point', coordinates: [0, 0]}"
	}

	return ""
}

func (s TypeData) Str() string {
	switch s {
	case TypeDataArray:
		return "array"
	case TypeDataBool:
		return "bool"
	case TypeDataInt:
		return "int"
	case TypeDataKey:
		return "key"
	case TypeDataState:
		return "state"
	case TypeDataMemo:
		return "memo"
	case TypeDataNumber:
		return "number"
	case TypeDataPrecision:
		return "precision"
	case TypeDataObject:
		return "object"
	case TypeDataSerie:
		return "serie"
	case TypeDataIndex:
		return "index"
	case TypeDataShortText:
		return "short_text"
	case TypeDataText:
		return "text"
	case TypeDataTime:
		return "time"
	case TypeDataBytes:
		return "bytes"
	case TypeDataGeometry:
		return "geometry"
	case TypeDataFullText:
		return "full_text"
	case TypeDataSelect:
		return "select"
	default:
		return "text"
	}
}

type ColumnField string

const (
	PRIMARYKEY = "id"
	PK         = PRIMARYKEY
	KEY        = PRIMARYKEY
	INDEX      = "index"
	PROJECT    = "project_id"
	CREATED_AT = "created_at"
	UPDATED_AT = "updated_at"
	STATUS     = "status_id"
	SYSID      = "jdbid"
	CREATED_TO = "created_to"
	UPDATED_TO = "updated_to"
	FULLTEXT   = "fulltext"
	HISTORYCAL = "historical"
	ASC        = true
)

var (
	IndexField      ColumnField = INDEX
	SourceField     ColumnField = "source"
	ProjectField    ColumnField = PROJECT
	CreatedAtField  ColumnField = CREATED_AT
	UpdatedAtField  ColumnField = UPDATED_AT
	StatusField     ColumnField = STATUS
	PrimaryKeyField ColumnField = PRIMARYKEY
	SystemKeyField  ColumnField = SYSID
	CreatedToField  ColumnField = CREATED_TO
	UpdatedToField  ColumnField = UPDATED_TO
	FullTextField   ColumnField = FULLTEXT
)

func (s ColumnField) TypeData() TypeData {
	switch s {
	case IndexField:
		return TypeDataIndex
	case SourceField:
		return TypeDataObject
	case ProjectField:
		return TypeDataKey
	case CreatedAtField:
		return TypeDataTime
	case UpdatedAtField:
		return TypeDataTime
	case StatusField:
		return TypeDataState
	case PrimaryKeyField:
		return TypeDataKey
	case SystemKeyField:
		return TypeDataKey
	case CreatedToField:
		return TypeDataTime
	case UpdatedToField:
		return TypeDataTime
	case FullTextField:
		return TypeDataFullText
	}

	return TypeDataText
}

func (s ColumnField) Str() string {
	return string(s)
}

type Relation struct {
	With            *Model            `json:"-"`
	Fk              map[string]string `json:"fk"`
	Limit           int               `json:"rows"`
	OnDeleteCascade bool              `json:"on_delete_cascade"`
	OnUpdateCascade bool              `json:"on_update_cascade"`
}

/**
* Where
* @param from et.Json
* @return et.Json
**/
func (s *Relation) Where(from et.Json) et.Json {
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

type Rollup struct {
	With   *Model            `json:"-"`
	Fk     map[string]string `json:"fk"`
	Fields []interface{}     `json:"fields"`
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

	return result
}

type FullText struct {
	Language string   `json:"language"`
	Columns  []string `json:"columns"`
}

func (s *FullText) describe() et.Json {
	result := et.Json{
		"language": s.Language,
		"columns":  s.Columns,
	}

	return result
}

type Column struct {
	Model        *Model                  `json:"-"`
	Source       *Column                 `json:"-"`
	Name         string                  `json:"name"`
	Description  string                  `json:"description"`
	TypeColumn   TypeColumn              `json:"type_column"`
	TypeData     TypeData                `json:"type_data"`
	Default      interface{}             `json:"default"`
	Max          float64                 `json:"max"`
	Min          float64                 `json:"min"`
	Hidden       bool                    `json:"hidden"`
	Detail       *Relation               `json:"detail"`
	Rollup       *Rollup                 `json:"rollup"`
	FullText     *FullText               `json:"fulltext"`
	Values       interface{}             `json:"values"`
	CalcFunction map[string]DataFunction `json:"-"`
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
		Model:        model,
		Name:         name,
		Description:  description,
		TypeColumn:   typeColumn,
		TypeData:     typeData,
		Default:      def,
		Max:          0,
		Min:          0,
		Values:       "",
		CalcFunction: make(map[string]DataFunction, 0),
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
func (s *Column) SetValue(value interface{}) {
	s.Values = value
}

/**
* DefaultValue
* @return interface{}
**/
func (s *Column) DefaultValue() interface{} {
	if s.TypeData == TypeDataTime {
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
	result = utility.Quote(result)

	return result
}

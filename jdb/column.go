package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/timezone"
	"github.com/cgalvisleon/et/utility"
)

type TypeColumn int

const (
	TpColumn TypeColumn = iota
	TpAtribute
	TpDetail
	TpGenerated
	TpRelation
	TpRollup
)

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
	TypeDataBool
	TypeDataTime
	TypeDataBytes
	// Special
	TypeDataObject
	TypeDataArray
	TypeDataGeometry
	TypeDataFullText
	TypeDataNone
	TypeListValues
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
	default:
		return "text"
	}
}

type ColumnField string

const (
	INDEX         = "index"
	SOURCE        = "source"
	PROJECT       = "project"
	CREATED_AT    = "created_at"
	UPDATED_AT    = "update_at"
	STATUS        = "status"
	KEY           = "id"
	SYSID         = "jdbid"
	CREATED_TO    = "created_to"
	UPDATED_TO    = "updated_to"
	FULLTEXT      = "fulltext"
	HISTORY_INDEX = "history_index"
)

var (
	IndexField     ColumnField = INDEX
	SourceField    ColumnField = SOURCE
	ProjectField   ColumnField = PROJECT
	CreatedAtField ColumnField = CREATED_AT
	UpdatedAtField ColumnField = UPDATED_AT
	StateField     ColumnField = STATUS
	KeyField       ColumnField = KEY
	SystemKeyField ColumnField = SYSID
	CreatedToField ColumnField = CREATED_TO
	UpdatedToField ColumnField = UPDATED_TO
	FullTextField  ColumnField = FULLTEXT
)

func (s ColumnField) TypeData() TypeData {
	switch s {
	case IndexField:
		return TypeDataInt
	case SourceField:
		return TypeDataObject
	case ProjectField:
		return TypeDataKey
	case CreatedAtField:
		return TypeDataTime
	case UpdatedAtField:
		return TypeDataTime
	case StateField:
		return TypeDataState
	case KeyField:
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

type Column struct {
	Model         *Model        `json:"-"`
	Name          string        `json:"name"`
	Field         string        `json:"field"`
	Description   string        `json:"description"`
	Table         string        `json:"table"`
	TypeColumn    TypeColumn    `json:"type_column"`
	TypeData      TypeData      `json:"type_data"`
	Default       interface{}   `json:"default"`
	Max           float64       `json:"max"`
	Min           float64       `json:"min"`
	Hidden        bool          `json:"hidden"`
	FullText      []string      `json:"columns"`
	Detail        *Detail       `json:"detail"`
	FuncGenerated FuncGenerated `json:"-"`
	Limit         int           `json:"limit"`
	Language      string        `json:"language"`
	Values        []interface{} `json:"values"`
}

func newColumn(model *Model, name string, description string, typeColumn TypeColumn, typeData TypeData, def interface{}) *Column {
	return &Column{
		Model:       model,
		Name:        name,
		Field:       name,
		Description: description,
		Table:       model.Table,
		TypeColumn:  typeColumn,
		TypeData:    typeData,
		Default:     def,
		Max:         0,
		Min:         0,
		Hidden:      false,
		FullText:    []string{},
		Limit:       30,
		Values:      []interface{}{},
	}
}

/**
* Describe
* @return et.Json
**/
func (s *Column) Describe() et.Json {
	result, err := et.Object(s)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* Fk
* @return string
**/
func (s *Column) Fk() string {
	result := strs.ReplaceAll(s.Field, []string{"_"}, "")
	result = s.Model.Name + "_" + result

	return result
}

/**
* DefaultValue
* @return interface{}
**/
func (s *Column) DefaultValue() interface{} {
	if s.Name == string(ProjectField) {
		return "-1"
	}
	switch s.TypeData {
	case TypeDataKey:
		return utility.UUID()
	case TypeDataTime:
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

/**
* GetField
* @return *Field
**/
func (s *Column) GetField() *Field {
	result := &Field{
		Column: s,
		Schema: s.Model.Schema.Name,
		Table:  s.Model.Name,
		Field:  s.Field,
		Name:   s.Name,
		Alias:  s.Name,
	}
	if s.TypeColumn != TpColumn {
		result.Atrib = s.Name
	}

	return result
}

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
	TpGenerate
	TpDetail
)

type TypeData int

const (
	TypeDataText TypeData = iota
	TypeDataMemo
	TypeDataShortText
	TypeDataKey
	TypeDataInt
	TypeDataNumber
	TypeDataPrecision
	TypeDataSerie
	TypeDataBool
	TypeDataTime
	// Special
	TypeDataObject
	TypeDataArray
	TypeDataFullText
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
	case TypeDataMemo:
		return ""
	case TypeDataNumber:
		return 0.0
	case TypeDataPrecision:
		return 0.0
	case TypeDataObject:
		return "{}"
	case TypeDataSerie:
		return 0
	case TypeDataShortText:
		return ""
	case TypeDataText:
		return ""
	case TypeDataTime:
		return "NOW()"
	}

	return ""
}

type ColumnField string

var (
	IndexField     ColumnField = "index"
	SourceField    ColumnField = "_data"
	ProjectField   ColumnField = "project_id"
	CreatedAtField ColumnField = "created_at"
	UpdatedAtField ColumnField = "update_at"
	StateField     ColumnField = "_state"
	KeyField       ColumnField = "_id"
	SystemKeyField ColumnField = "_idt"
	ClassField     ColumnField = "_class"
	CreatedToField ColumnField = "created_to"
	UpdatedToField ColumnField = "updated_to"
	FullTextField  ColumnField = "_fulltext"
)

func (s ColumnField) Str() string {
	return string(s)
}

func (s ColumnField) Up() string {
	return strs.Uppcase(string(s))
}

func (s ColumnField) Low() string {
	return strs.Lowcase(string(s))
}

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
		return TypeDataKey
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
	Model       *Model      `json:"-"`
	Name        string      `json:"name"`
	Field       string      `json:"field"`
	Description string      `json:"description"`
	Table       string      `json:"table"`
	TypeColumn  TypeColumn  `json:"type_column"`
	TypeData    TypeData    `json:"type_data"`
	Default     interface{} `json:"default"`
	Max         float64     `json:"max"`
	Min         float64     `json:"min"`
	Hidden      bool        `json:"hidden"`
	Columns     []string    `json:"columns"`
	Definition  interface{} `json:"definition"`
}

func newColumn(model *Model, name string, description string, typeColumn TypeColumn, typeData TypeData, def interface{}) *Column {
	return &Column{
		Model:       model,
		Name:        Name(name),
		Field:       strs.Uppcase(name),
		Description: description,
		Table:       model.Table,
		TypeColumn:  typeColumn,
		TypeData:    typeData,
		Default:     def,
		Max:         0,
		Min:         0,
		Hidden:      false,
		Columns:     []string{},
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
	return s.Model.Table + "_" + s.Field
}

/**
* Up
* @return string
**/
func (s *Column) Up() string {
	return strs.Uppcase(s.Field)
}

/**
* Low
* @return string
**/
func (s *Column) Low() string {
	return strs.Lowcase(s.Field)
}

/**
* DefaultValue
* @return interface{}
**/
func (s *Column) DefaultValue() interface{} {
	switch s.TypeData {
	case TypeDataKey:
		return utility.UUID()
	case TypeDataTime:
		return timezone.Now()
	}

	return s.Default
}

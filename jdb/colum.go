package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

func fieldName(name string) string {
	return strs.Uppcase(name)
}

type TypeColumn int

const (
	TpColumn TypeColumn = iota
	TpAtribute
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
	TypeFullText
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
	case TypeFullText:
		return ""
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

func (s ColumnField) Uppcase() string {
	return strs.Uppcase(string(s))
}

func (s ColumnField) Lowcase() string {
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
		return TypeFullText
	}

	return TypeDataText
}

type Column struct {
	Model       *Model
	Name        string
	Description string
	Table       string
	Field       string
	TypeColumn  TypeColumn
	TypeData    TypeData
	Default     interface{}
	Max         float64
	Min         float64
	Hidden      bool
	Columns     []string
}

func newColumn(model *Model, name string, description string, typeColumn TypeColumn, typeData TypeData, def interface{}) *Column {
	return &Column{
		Model:       model,
		Name:        name,
		Description: description,
		Table:       model.Table,
		Field:       fieldName(name),
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

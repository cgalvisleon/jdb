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
	TypeDataSerie
	TypeDataBool
	TypeDataTime
	// Special
	TypeDataObject
	TypeDataArray
	TypeFullText
)

func (s TypeData) DefaultValue(d Driver) interface{} {
	return d.DefaultValue(s)
}

type ColumnField string

var (
	IndexField     ColumnField = "index"
	DataField      ColumnField = "_data"
	ProjectField   ColumnField = "project_id"
	CreatedAtField ColumnField = "created_at"
	UpdatedAtField ColumnField = "update_at"
	StateField     ColumnField = "_state"
	KeyField       ColumnField = "_id"
	SystemKeyField ColumnField = "_idt"
	CreatedToField ColumnField = "created_to"
	UpdatedToField ColumnField = "updated_to"
	FullTextField  ColumnField = "_fulltext"
)

func (s ColumnField) Str() string {
	return string(s)
}

func (s ColumnField) TypeData() TypeData {
	switch s {
	case IndexField:
		return TypeDataInt
	case DataField:
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

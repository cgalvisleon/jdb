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
	TpDetails
)

type TypeData int

const (
	TypeDataText TypeData = iota
	TypeDataMemo
	TypeDataShort
	TypeDataKey
	TypeDataInt
	TypeDataFloat
	TypeDataNumber
	TypeDataSerie
	TypeDataBool
	TypeDataTime
	// Special
	TypeDataSource
	TypeDataObject
	TypeDataArray
	TypeFullText
)

type TypeDefault int

const (
	DefaultNow TypeDefault = iota
	DefaultToday
	DefaultTime
	DefaultNone
	DefaultZero
	DefaultId
	DefaultKey
)

var (
	IndexField     = "index"
	DataField      = "_data"
	ProjectField   = "project_id"
	CreatedAtField = "created_at"
	UpdatedAtField = "update_at"
	StateField     = "_state"
	SystemKeyField = "_idt"
	CreatedToField = "created_to"
	UpdatedToField = "updated_to"
	FullTextField  = "_fulltext"
)

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

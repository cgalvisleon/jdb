package jdb

import "github.com/cgalvisleon/et/strs"

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
	TypeDataString TypeData = iota
	TypeDataInt
	TypeDataFloat
	TypeDataNumber
	TypeDataBool
	TypeDataTime
	TypeDataText
	// Special
	TypeDataSource
	TypeDataObject
	TypeDataArray
)

type Column struct {
	Model      *Model
	Name       string
	Describe   string
	Table      string
	Field      string
	TypeColumn TypeColumn
	TypeData   TypeData
	Default    interface{}
	Max        float64
	Min        float64
}

func NewColumn(model *Model, name string, describe string, typeColumn TypeColumn, typeData TypeData, def interface{}) *Column {
	return &Column{
		Model:      model,
		Name:       name,
		Describe:   describe,
		Table:      model.Table,
		Field:      fieldName(name),
		TypeColumn: typeColumn,
		TypeData:   typeData,
		Default:    def,
		Max:        0,
		Min:        0,
	}
}

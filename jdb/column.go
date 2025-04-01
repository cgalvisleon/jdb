package jdb

import (
	"encoding/gob"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/timezone"
	"github.com/cgalvisleon/et/utility"
)

type TypeColumn int

const (
	TpColumn TypeColumn = iota
	TpAtribute
	TpGenerated
	TpRelatedTo
	TpRollup
	TpIA
)

func (s TypeColumn) Str() string {
	switch s {
	case TpColumn:
		return "column"
	case TpAtribute:
		return "attribute"
	case TpGenerated:
		return "generated"
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
	TypeDataBool
	TypeDataTime
	TypeDataBytes
	// Special
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
	case TypeDataSelect:
		return "select"
	default:
		return "text"
	}
}

type ColumnField string

const (
	PRIMARYKEY    = "id"
	PK            = PRIMARYKEY
	KEY           = PRIMARYKEY
	SOURCE        = "source"
	INDEX         = "index"
	PROJECT       = "project_id"
	CREATED_AT    = "created_at"
	UPDATED_AT    = "update_at"
	STATUS        = "status"
	SYSID         = "jdbid"
	CREATED_TO    = "created_to"
	UPDATED_TO    = "updated_to"
	FULLTEXT      = "fulltext"
	HISTORY_INDEX = "hindex"
)

var (
	IndexField      ColumnField = INDEX
	SourceField     ColumnField = SOURCE
	ProjectField    ColumnField = PROJECT
	CreatedAtField  ColumnField = CREATED_AT
	UpdatedAtField  ColumnField = UPDATED_AT
	StateField      ColumnField = STATUS
	PrimaryKeyField ColumnField = PRIMARYKEY
	SystemKeyField  ColumnField = SYSID
	CreatedToField  ColumnField = CREATED_TO
	UpdatedToField  ColumnField = UPDATED_TO
	FullTextField   ColumnField = FULLTEXT
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

type FullText struct {
	Language string    `json:"language"`
	Columns  []*Column `json:"columns"`
}

type GeneratedFunction func(col *Column, data *et.Json)

type Relation struct {
	Key             string  `json:"key"`
	With            *Model  `json:"with"`
	Fk              *Column `json:"fk"`
	Limit           int     `json:"rows"`
	OnDeleteCascade bool    `json:"on_delete_cascade"`
	OnUpdateCascade bool    `json:"on_update_cascade"`
}

func (s *Relation) Describe() et.Json {
	with := ""
	if s.With != nil {
		with = s.With.Name
	}
	fk := ""
	if s.Fk != nil {
		fk = s.Fk.Name
	}
	result := et.Json{
		"key":               s.Key,
		"with":              with,
		"fk":                fk,
		"rows":              s.Limit,
		"on_delete_cascade": s.OnDeleteCascade,
		"on_update_cascade": s.OnUpdateCascade,
	}

	return result
}

type Rollup struct {
	Key    string  `json:"key"`
	Source *Model  `json:"source"`
	Fk     *Column `json:"fk"`
	Props  []*Column
}

func (s *Rollup) Describe() et.Json {
	source := ""
	if s.Source != nil {
		source = s.Source.Name
	}
	props := make([]string, 0)
	for _, column := range s.Props {
		props = append(props, column.Name)
	}

	result := et.Json{
		"key":    s.Key,
		"source": source,
		"fk":     s.Fk.Name,
		"props":  props,
	}

	return result
}

type Column struct {
	Model             *Model            `json:"-"`
	Source            *Column           `json:"-"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	TypeColumn        TypeColumn        `json:"type_column"`
	TypeData          TypeData          `json:"type_data"`
	Default           interface{}       `json:"default"`
	Max               float64           `json:"max"`
	Min               float64           `json:"min"`
	Hidden            bool              `json:"hidden"`
	Detail            *Relation         `json:"detail"`
	Rollup            *Rollup           `json:"rollup"`
	FullText          *FullText         `json:"columns"`
	GeneratedFunction GeneratedFunction `json:"-"`
	Values            []interface{}     `json:"values"`
}

func init() {
	gob.Register(&Column{})
}

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
		Values:      []interface{}{},
	}
}

/**
* Describe
* @return et.Json
**/
func (s *Column) Describe() et.Json {
	var fulltext = []string{}
	if s.FullText != nil {
		for _, column := range s.FullText.Columns {
			fulltext = append(fulltext, column.Name)
		}
	}
	detail := et.Json{}
	if s.Detail != nil {
		detail = s.Detail.Describe()
	}
	rollup := et.Json{}
	if s.Rollup != nil {
		rollup = s.Rollup.Describe()
	}

	result := et.Json{
		"name":        s.Name,
		"description": s.Description,
		"type_column": s.TypeColumn.Str(),
		"type_data":   s.TypeData.DefaultValue(),
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
func (s *Column) Idx() int {
	if s.Model == nil {
		return -1
	}

	return slices.IndexFunc(s.Model.Columns, func(e *Column) bool { return e == s })
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

package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/et"
)

/**
* setSelect
* @param field *Field
* @return *Ql
**/
func (s *Ql) setSelect(field *Field) *Ql {
	if field != nil {
		if field.Column == nil {
			return s
		}

		if slices.Contains([]TypeColumn{TpColumn, TpAtribute}, field.Column.TypeColumn) {
			s.Selects = append(s.Selects, field)
		} else if slices.Contains([]TypeColumn{TpRelatedTo}, field.Column.TypeColumn) {
			s.Details = append(s.Details, field)
		}
	}
	return s
}

/**
* Select
* @param fields ...string
* @return *Ql
**/
func (s *Ql) Select(fields ...string) *Ql {
	s.TypeSelect = Select
	for _, field := range fields {
		field := s.getField(field)
		s.setSelect(field)
	}

	return s
}

/**
* Data
* @param fields ...string
* @return *Ql
**/
func (s *Ql) Data(fields ...string) *Ql {
	result := s.Select(fields...)
	result.TypeSelect = Data

	return result
}

/**
* Exec
* @return et.Items, error
**/
func (s *Ql) Exec() (et.Items, error) {
	return et.Items{}, nil
}

/**
* ListSelects
* @return []string
**/
func (s *Ql) listSelects() []string {
	result := []string{}
	for _, sel := range s.Selects {
		result = append(result, sel.AsField())
	}

	return result
}

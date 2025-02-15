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
		s.Selects = append(s.Selects, field)
		if field.Column != nil && !slices.Contains([]TypeColumn{}, field.Column.TypeColumn) {
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

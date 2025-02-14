package jdb

import (
	"github.com/cgalvisleon/et/et"
)

/**
* Select
* @param fields ...string
* @return *Ql
**/
func (s *Ql) Select(fields ...string) *Ql {
	for _, field := range fields {
		field := s.GetField(field)
		if field != nil {
			s.Selects = append(s.Selects, field)
		}
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

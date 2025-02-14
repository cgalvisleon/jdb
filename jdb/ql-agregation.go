package jdb

import (
	"github.com/cgalvisleon/et/strs"
)

/**
* Sum
* @param field string
* @return *Ql
**/
func (s *Ql) Sum(field string) *Ql {
	sel := s.GetSelect(field)
	if sel != nil {
		sel.Field.Agregation = AgregationSum
		sel.Field.Alias = strs.Format(`sum_%s`, sel.Field.Name)
	}

	return s
}

/**
* Count
* @param field string
* @return *Ql
**/
func (s *Ql) Count(field string) *Ql {
	sel := s.GetSelect(field)
	if sel != nil {
		sel.Field.Agregation = AgregationCount
		sel.Field.Alias = strs.Format(`count_%s`, sel.Field.Name)
	}

	return s
}

/**
* Avg
* @param field string
* @return *Ql
**/
func (s *Ql) Avg(field string) *Ql {
	sel := s.GetSelect(field)
	if sel != nil {
		sel.Field.Agregation = AgregationAvg
		sel.Field.Alias = strs.Format(`avg_%s`, sel.Field.Name)
	}

	return s
}

/**
* Min
* @param field string
* @return *Ql
**/
func (s *Ql) Min(field string) *Ql {
	sel := s.GetSelect(field)
	if sel != nil {
		sel.Field.Agregation = AgregationMin
		sel.Field.Alias = strs.Format(`min_%s`, sel.Field.Name)
	}

	return s
}

/**
* Max
* @param field string
* @return *Ql
**/
func (s *Ql) Max(field string) *Ql {
	sel := s.GetSelect(field)
	if sel != nil {
		sel.Field.Agregation = AgregationMax
		sel.Field.Alias = strs.Format(`max_%s`, sel.Field.Name)
	}

	return s
}

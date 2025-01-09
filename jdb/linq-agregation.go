package jdb

import (
	"regexp"

	"github.com/cgalvisleon/et/strs"
)

type TypeAgregation int

const (
	Nag TypeAgregation = iota
	AgregationSum
	AgregationCount
	AgregationAvg
	AgregationMin
	AgregationMax
)

type Agregation struct {
	Agregation string
	pattern    string
	re         *regexp.Regexp
}

var agregations = map[TypeAgregation]*Agregation{
	Nag:             {Agregation: "", pattern: ""},
	AgregationSum:   {Agregation: "SUM", pattern: `SUM\([a-zA-Z0-9_]+\)$`},
	AgregationCount: {Agregation: "COUNT", pattern: `COUNT\([a-zA-Z0-9_]+\)$`},
	AgregationAvg:   {Agregation: "AVG", pattern: `AVG\([a-zA-Z0-9_]+\)$`},
	AgregationMin:   {Agregation: "MIN", pattern: `MIN\([a-zA-Z0-9_]+\)$`},
	AgregationMax:   {Agregation: "MAX", pattern: `MAX\([a-zA-Z0-9_]+\)$`},
}

/**
* init
**/
func init() {
	for _, agregation := range agregations {
		re, err := regexp.Compile(agregation.pattern)
		if err != nil {
			continue
		}
		agregation.re = re
	}
}

/**
* Sum
* @param field string
* @return *Linq
**/
func (s *Linq) Sum(field string) *Linq {
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
* @return *Linq
**/
func (s *Linq) Count(field string) *Linq {
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
* @return *Linq
**/
func (s *Linq) Avg(field string) *Linq {
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
* @return *Linq
**/
func (s *Linq) Min(field string) *Linq {
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
* @return *Linq
**/
func (s *Linq) Max(field string) *Linq {
	sel := s.GetSelect(field)
	if sel != nil {
		sel.Field.Agregation = AgregationMax
		sel.Field.Alias = strs.Format(`max_%s`, sel.Field.Name)
	}

	return s
}

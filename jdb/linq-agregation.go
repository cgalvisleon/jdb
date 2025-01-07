package jdb

import "regexp"

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
	}

	return s
}

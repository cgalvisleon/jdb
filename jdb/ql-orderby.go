package jdb

import (
	"github.com/cgalvisleon/et/et"
)

type QlOrder struct {
	Asc  []*Field
	Desc []*Field
}

/**
* OrderBy
* @param sorted bool
* @param columns ...string
* @return *Ql
**/
func (s *Ql) OrderBy(sorted bool, columns ...string) *Ql {
	for _, col := range columns {
		field := s.getField(col, false)
		if field != nil {
			if sorted {
				s.Orders.Asc = append(s.Orders.Asc, field)
			} else {
				s.Orders.Desc = append(s.Orders.Desc, field)
			}
		}
	}

	return s
}

/**
* OrderByAsc
* @param columns ...any
* @return *Ql
**/
func (s *Ql) OrderByAsc(columns ...string) *Ql {
	return s.OrderBy(true, columns...)
}

/**
* OrderByDesc
* @param columns ...any
* @return *Ql
**/
func (s *Ql) OrderByDesc(columns ...string) *Ql {
	return s.OrderBy(false, columns...)
}

/**
* setOrders
* @param orders []et.Json
* @return *Ql
**/
func (s *Ql) setOrders(orders et.Json) *Ql {
	for key := range orders {
		switch key {
		case "asc", "ASC":
			val := orders.ArrayStr(key)
			s.OrderByAsc(val...)
		case "desc", "DESC":
			val := orders.ArrayStr(key)
			s.OrderByDesc(val...)
		}
	}

	return s
}

/**
* listOrders
* @return []string
**/
func (s *Ql) listOrders() et.Json {
	asc := []string{}
	desc := []string{}
	for _, sel := range s.Orders.Asc {
		asc = append(asc, sel.AsField())
	}
	for _, sel := range s.Orders.Desc {
		desc = append(desc, sel.AsField())
	}

	return et.Json{
		"asc":  asc,
		"desc": desc,
	}
}

package jdb

import (
	"github.com/cgalvisleon/et/et"
)

type QlOrder struct {
	Asc  []*Field
	Desc []*Field
}

/**
* orderBy
* @param asc bool, columns ...string
* @return *Ql
**/
func (s *Ql) orderBy(asc bool, columns ...string) *Ql {
	for _, col := range columns {
		field := s.getField(col)
		if field != nil {
			if asc {
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
	return s.orderBy(true, columns...)
}

/**
* OrderByDesc
* @param columns ...any
* @return *Ql
**/
func (s *Ql) OrderByDesc(columns ...string) *Ql {
	return s.orderBy(false, columns...)
}

/**
* OrderBy
* @param columns ...any
* @return *Ql
**/
func (s *Ql) OrderBy(columns ...string) *Ql {
	return s.OrderByAsc(columns...)
}

/**
* setOrderBy
* @param orders et.Json
* @return *Ql
**/
func (s *Ql) setOrderBy(orders et.Json) *Ql {
	if len(orders) == 0 {
		return s
	}

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
* getOrders
* @return []string
**/
func (s *Ql) getOrders() et.Json {
	asc := []string{}
	desc := []string{}
	for _, sel := range s.Orders.Asc {
		asc = append(asc, sel.asField())
	}
	for _, sel := range s.Orders.Desc {
		desc = append(desc, sel.asField())
	}

	return et.Json{
		"asc":  asc,
		"desc": desc,
	}
}

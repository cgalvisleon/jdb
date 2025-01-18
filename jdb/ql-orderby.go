package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type QlOrder struct {
	Field  *Field
	Sorted bool
}

/**
* OrderBy
* @param sorted bool
* @param columns ...string
* @return *Ql
**/
func (s *Ql) OrderBy(sorted bool, columns ...string) *Ql {
	for _, col := range columns {
		field := s.GetField(col, true)
		if field != nil {
			order := &QlOrder{
				Field:  field,
				Sorted: sorted,
			}
			s.Orders = append(s.Orders, order)
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
func (s *Ql) setOrders(orders []et.Json) *Ql {
	for _, item := range orders {
		sorted := item.Bool("sorted")
		columns := item.ArrayStr("columns")
		s.OrderBy(sorted, columns...)
	}

	return s
}

/**
* listOrders
* @return []string
**/
func (s *Ql) listOrders() []string {
	result := []string{}
	for _, sel := range s.Orders {
		result = append(result, strs.Format(`%s, SORTED:%v`, sel.Field.AsField(), sel.Sorted))
	}

	return result
}

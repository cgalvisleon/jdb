package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type LinqOrder struct {
	Field  *Field
	Sorted bool
}

/**
* OrderBy
* @param sorted bool
* @param columns ...string
* @return *Linq
**/
func (s *Linq) OrderBy(sorted bool, columns ...string) *Linq {
	for _, col := range columns {
		field := s.GetField(col, true)
		if field != nil {
			order := &LinqOrder{
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
* @return *Linq
**/
func (s *Linq) OrderByAsc(columns ...string) *Linq {
	return s.OrderBy(true, columns...)
}

/**
* OrderByDesc
* @param columns ...any
* @return *Linq
**/
func (s *Linq) OrderByDesc(columns ...string) *Linq {
	return s.OrderBy(false, columns...)
}

/**
* setOrders
* @param orders []et.Json
* @return *Linq
**/
func (s *Linq) setOrders(orders []et.Json) *Linq {
	for _, item := range orders {
		sorted := item.Bool("sorted")
		columns := item.ArrayStr([]string{}, "columns")
		s.OrderBy(sorted, columns...)
	}

	return s
}

/**
* listOrders
* @return []string
**/
func (s *Linq) listOrders() []string {
	result := []string{}
	for _, sel := range s.Orders {
		result = append(result, strs.Format(`%s, SORTED:%v`, sel.Field.AsField(), sel.Sorted))
	}

	return result
}

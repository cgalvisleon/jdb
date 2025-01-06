package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type LinqOrder struct {
	LinqSelect
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
		c := s.GetSelect(col)
		if c != nil {
			order := &LinqOrder{
				LinqSelect: *c,
				Sorted:     sorted,
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
* SetOrders
* @param orders []et.Json
* @return *Linq
**/
func (s *Linq) SetOrders(orders []et.Json) *Linq {
	for _, item := range orders {
		sorted := item.Bool("sorted")
		columns := item.ArrayStr([]string{}, "columns")
		s.OrderBy(sorted, columns...)
	}

	return s
}

/**
* ListOrders
* @return []string
**/
func (s *Linq) ListOrders() []string {
	result := []string{}
	for _, sel := range s.Orders {
		result = append(result, strs.Format(`%s, SORTED:%v`, sel.Field.AsField(), sel.Sorted))
	}

	return result
}

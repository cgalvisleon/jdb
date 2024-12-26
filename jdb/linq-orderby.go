package jdb

import "github.com/cgalvisleon/et/strs"

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
		c := s.getSelect(col)
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

func (s *Linq) ListOrders() []string {
	result := []string{}
	for _, sel := range s.Orders {
		result = append(result, strs.Format(`%s, %s, SORTED:%v`, sel.Field.Tag(), sel.Field.Caption(), sel.Sorted))
	}

	return result
}

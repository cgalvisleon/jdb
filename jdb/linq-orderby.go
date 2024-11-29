package jdb

/**
* OrderBy
* @param sorted bool
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) OrderBy(sorted bool, columns ...interface{}) *Linq {
	for _, col := range columns {
		c := s.getColumn(col)
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
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) OrderByAsc(columns ...interface{}) *Linq {
	return s.OrderBy(true, columns...)
}

/**
* OrderByDesc
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) OrderByDesc(columns ...interface{}) *Linq {
	return s.OrderBy(false, columns...)
}

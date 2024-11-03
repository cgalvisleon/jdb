package jdb

/**
* OrderByAsc
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) OrderByAsc(columns ...interface{}) *Linq {
	for _, col := range columns {
		c := s.getColumn(col)
		if c != nil {
			order := &LinqOrder{
				LinqSelect: *c,
				Sorted:     true,
			}
			s.Orders = append(s.Orders, order)
		}
	}

	return s
}

/**
* OrderByDesc
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) OrderByDesc(columns ...interface{}) *Linq {
	for _, col := range columns {
		c := s.getColumn(col)
		if c != nil {
			order := &LinqOrder{
				LinqSelect: *c,
				Sorted:     false,
			}
			s.Orders = append(s.Orders, order)
		}
	}

	return s
}

package jdb

/**
* Select
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) Select(columns ...interface{}) *Linq {
	s.TypeLinq = TypeLinqSelect
	for _, col := range columns {
		c := s.getColumn(col)
		if c != nil {
			s.Selects = append(s.Selects, c)
		}
	}

	return s
}

/**
* Data
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) Data(columns ...interface{}) *Linq {
	s.TypeLinq = TypeLinqData
	for _, col := range columns {
		c := s.getColumn(col)
		if c != nil {
			s.Selects = append(s.Selects, c)
		}
	}

	return s
}

/**
* Offset
* @param offset int
* @return *Linq
**/
func (s *Linq) Page(val int) *Linq {
	s.page = val
	s.Offset = s.Limit * (s.page - 1)
	return s
}

/**
* Limit
* @param limit int
* @return *Linq
**/
func (s *Linq) Rows(val int) *Linq {
	s.Limit = val
	s.Offset = s.Limit * (s.page - 1)
	return s
}

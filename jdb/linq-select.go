package jdb

/**
* Select
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) Select(columns ...interface{}) *Linq {
	s.TypeLinq = TypeLinqSelect
	for _, column := range columns {
		col := s.getColumn(column)
		if col != nil {
			s.Selects = append(s.Selects, col)
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
	for _, column := range columns {
		col := s.getColumn(column)
		if col != nil {
			s.Selects = append(s.Selects, col)
		}
	}

	return s
}

/**
* Sum
* @param column interface{}
* @return *Linq
**/
func (s *Linq) Sum(column interface{}) *Linq {
	col := s.getColumn(column)
	col.Function = Sum
	if col != nil {
		s.Selects = append(s.Selects, col)
	}

	return s
}

/**
* Count
* @param column interface{}
* @return *Linq
**/
func (s *Linq) Count(column interface{}) *Linq {
	col := s.getColumn(column)
	col.Function = Count
	if col != nil {
		s.Selects = append(s.Selects, col)
	}

	return s
}

/**
* Avg
* @param column interface{}
* @return *Linq
**/
func (s *Linq) Avg(column interface{}) *Linq {
	col := s.getColumn(column)
	col.Function = Avg
	if col != nil {
		s.Selects = append(s.Selects, col)
	}

	return s
}

/**
* Max
* @param column interface{}
* @return *Linq
**/
func (s *Linq) Max(column interface{}) *Linq {
	col := s.getColumn(column)
	col.Function = Max
	if col != nil {
		s.Selects = append(s.Selects, col)
	}

	return s
}

/**
* Min
* @param column interface{}
* @return *Linq
**/
func (s *Linq) Min(column interface{}) *Linq {
	col := s.getColumn(column)
	col.Function = Min
	if col != nil {
		s.Selects = append(s.Selects, col)
	}

	return s
}

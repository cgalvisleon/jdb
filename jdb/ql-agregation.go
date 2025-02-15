package jdb

/**
* Sum
* @param field string
* @return *Ql
**/
func (s *Ql) Sum(field string) *Ql {
	agr := s.getField(field)
	if agr != nil {
		agr.SetAgregation(AgregationSum)
		s.setSelect(agr)
	}

	return s
}

/**
* Count
* @param field string
* @return *Ql
**/
func (s *Ql) Count(field string) *Ql {
	agr := s.getField(field)
	if agr != nil {
		agr.SetAgregation(AgregationCount)
		s.setSelect(agr)
	}

	return s
}

/**
* Avg
* @param field string
* @return *Ql
**/
func (s *Ql) Avg(field string) *Ql {
	agr := s.getField(field)
	if agr != nil {
		agr.SetAgregation(AgregationAvg)
		s.setSelect(agr)
	}

	return s
}

/**
* Min
* @param field string
* @return *Ql
**/
func (s *Ql) Min(field string) *Ql {
	agr := s.getField(field)
	if agr != nil {
		agr.SetAgregation(AgregationMin)
		s.setSelect(agr)
	}

	return s
}

/**
* Max
* @param field string
* @return *Ql
**/
func (s *Ql) Max(field string) *Ql {
	agr := s.getField(field)
	if agr != nil {
		agr.SetAgregation(AgregationMax)
		s.setSelect(agr)
	}

	return s
}

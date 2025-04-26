package jdb

import "github.com/cgalvisleon/et/utility"

/**
* Sum
* @param field string
* @return *Ql
**/
func (s *Ql) Sum(field string) *Ql {
	agr := s.getField(field, false)
	if agr != nil {
		agr.setAgregation(AgregationSum)
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
	if !utility.ValidWord(field) {
		return s
	}

	agr := s.getField(field, false)
	if agr != nil {
		agr.setAgregation(AgregationCount)
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
	if !utility.ValidWord(field) {
		return s
	}

	agr := s.getField(field, false)
	if agr != nil {
		agr.setAgregation(AgregationAvg)
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
	if !utility.ValidWord(field) {
		return s
	}

	agr := s.getField(field, false)
	if agr != nil {
		agr.setAgregation(AgregationMin)
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
	if !utility.ValidWord(field) {
		return s
	}

	agr := s.getField(field, false)
	if agr != nil {
		agr.setAgregation(AgregationMax)
		s.setSelect(agr)
	}

	return s
}

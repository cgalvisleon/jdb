package jdb

import "github.com/cgalvisleon/et/utility"

/**
* Sum
* @param field string
* @return *Ql
**/
func (s *Ql) Sum(field string) *Ql {
	agr := s.getField(field)
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

	agr := s.getField(field)
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

	agr := s.getField(field)
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

	agr := s.getField(field)
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

	agr := s.getField(field)
	if agr != nil {
		agr.setAgregation(AgregationMax)
		s.setSelect(agr)
	}

	return s
}

/**
* Extract
* @param field string
* @return *Ql
**/
func (s *Ql) ExtractYear(field string) *Ql {
	if !utility.ValidWord(field) {
		return s
	}

	agr := s.getField(field)
	if agr != nil {
		agr.setAgregation(ExtractYear)
		s.setSelect(agr)
	}

	return s
}

/**
* ExtractMonth
* @param field string
* @return *Ql
**/
func (s *Ql) ExtractMonth(field string) *Ql {
	if !utility.ValidWord(field) {
		return s
	}

	agr := s.getField(field)
	if agr != nil {
		agr.setAgregation(ExtractMonth)
		s.setSelect(agr)
	}

	return s
}

/**
* ExtractDay
* @param field string
* @return *Ql
**/
func (s *Ql) ExtractDay(field string) *Ql {
	if !utility.ValidWord(field) {
		return s
	}

	agr := s.getField(field)
	if agr != nil {
		agr.setAgregation(ExtractDay)
		s.setSelect(agr)
	}

	return s
}

/**
* ExtractHour
* @param field string
* @return *Ql
**/
func (s *Ql) ExtractHour(field string) *Ql {
	if !utility.ValidWord(field) {
		return s
	}

	agr := s.getField(field)
	if agr != nil {
		agr.setAgregation(ExtractHour)
		s.setSelect(agr)
	}

	return s
}

/**
* ExtractMinute
* @param field string
* @return *Ql
**/
func (s *Ql) ExtractMinute(field string) *Ql {
	if !utility.ValidWord(field) {
		return s
	}

	agr := s.getField(field)
	if agr != nil {
		agr.setAgregation(ExtractMinute)
		s.setSelect(agr)
	}

	return s
}

/**
* ExtractSecond
* @param field string
* @return *Ql
**/
func (s *Ql) ExtractSecond(field string) *Ql {
	if !utility.ValidWord(field) {
		return s
	}

	agr := s.getField(field)
	if agr != nil {
		agr.setAgregation(ExtractSecond)
		s.setSelect(agr)
	}

	return s
}

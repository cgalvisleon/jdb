package jdb

import "github.com/cgalvisleon/et/et"

/**
* Select
* @param fields ...interface{}
* @return *Ql
**/
func (s *Model) Select(fields ...interface{}) *Ql {
	result := From(s)
	result.Select(fields...)

	return result
}

/**
* Data
* @param fields ...interface{}
* @return *Ql
**/
func (s *Model) Data(fields ...interface{}) *Ql {
	result := From(s)
	result.Data(fields...)

	return result
}

/**
* Where
* @param val string
* @return *Ql
**/
func (s *Model) Where(val string) *Ql {
	result := From(s)
	if s.SourceField != nil {
		result.TypeSelect = Source
	}

	return result.Where(val)
}

/**
* Join
* @param name string
* @return *Model
**/
func (s *Model) Join(name string) *QlJoin {
	return From(s).Join(name)
}

/**
* CountedTx
* @return int, error
**/
func (s *Model) CountedTx(tx *Tx) (int, error) {
	all, err := From(s).
		CountedTx(tx)
	if err != nil {
		return 0, err
	}

	return all, nil
}

/**
* Counted
* @return int, error
**/
func (s *Model) Counted() (int, error) {
	return s.CountedTx(nil)
}

/**
* QTx
* @param params et.Json
* @return et.Json, error
**/
func (s *Model) QTx(tx *Tx, params et.Json) (map[string]interface{}, error) {
	result, err := From(s).QueryTx(tx, params)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return result.ToMap(), nil
}

/**
* Q
* @param params et.Json
* @return map[string]interface{}, error
**/
func (s *Model) Q(params et.Json) (map[string]interface{}, error) {
	return s.QTx(nil, params)
}

package jdb

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
* GetFrom
* @return *QlFrom
**/
func (s *Model) GetFrom() *QlFrom {
	return &QlFrom{Model: s}
}

package jdb

/**
* Where
* @param val interface{}
* @return *Ql
**/
func (s *Model) Where(val interface{}) *Ql {
	result := From(s)
	if s.SourceField != nil {
		result.TypeSelect = Data
	}

	return result.Where(val)
}

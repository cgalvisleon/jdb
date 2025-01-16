package jdb

/**
* Where
* @param val interface{}
* @return *Linq
**/
func (s *Model) Where(val interface{}) *Linq {
	result := From(s)
	if s.SourceField != nil {
		result.TypeSelect = Data
	}

	return result.Where(val)
}

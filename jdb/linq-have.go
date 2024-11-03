package jdb

/**
* Having
* @param col interface{}
* @return *LinqWhere
**/
func (s *Linq) Having(col interface{}) *LinqWhere {
	result := &LinqWhere{
		Linq: s,
	}

	return result
}

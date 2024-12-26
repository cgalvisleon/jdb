package jdb

/**
* Having
* @param col ...string
* @return *LinqWhere
**/
func (s *Linq) Having(col ...string) *LinqWhere {
	result := &LinqWhere{}

	return result
}

func (s *Linq) ListHavings() []string {
	result := []string{}
	for _, val := range s.Havings {
		result = append(result, val.String())
	}

	return result
}

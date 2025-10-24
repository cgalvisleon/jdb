package jdb

import "github.com/cgalvisleon/et/et"

/**
* Having
* @param having et.Json
* @return *Ql
**/
func (s *Ql) Having(having ...Condition) *Ql {
	havings := make([]et.Json, 0)
	for _, v := range having {
		havings = append(havings, v.ToJson())
	}

	s.Havings = havings
	return s
}

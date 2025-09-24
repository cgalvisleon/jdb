package jdb

import "github.com/cgalvisleon/et/et"

/**
* Join
* @param to, as string, on []et.Json
* @return *Ql
**/
func (s *Ql) Join(to, as string, on []et.Json) *Ql {
	n := len(s.Froms)
	if n == 0 {
		return s
	}

	model, err := s.db.getModelByName(to)
	if err != nil {
		return s
	}

	n = len(s.Joins) + 1
	s.Joins = append(s.Joins, et.Json{
		"from": et.Json{
			model.Table: as,
		},
		"on": on,
	})

	return s
}

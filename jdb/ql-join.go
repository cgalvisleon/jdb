package jdb

import "github.com/cgalvisleon/et/et"

/**
* Join
* @param to *Model, as string, on Condition
* @return *Ql
*
 */
func (s *Ql) Join(to *Model, as string, on *Condition) *Ql {
	n := len(s.Froms)
	if n == 0 {
		return s
	}

	n = len(s.Joins) + 1
	s.Joins = append(s.Joins, et.Json{
		"from": et.Json{
			to.Table: as,
		},
		"on": []et.Json{
			on.ToJson(),
		},
	})
	s.useJoin = true

	return s
}

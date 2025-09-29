package jdb

/**
* Where
* @param cond Condition
* @return *Cmd
**/
func (s *Ql) Where(cond Condition) *Ql {
	s.where.Where(cond)
	s.useJoin = false
	return s
}

/**
* whereJoin
* @param cond Condition
* @param conector string
* @return *Cmd
**/
func (s *Ql) whereJoin(cond Condition, conector string) *Ql {
	n := len(s.Joins) - 1
	if n < 0 {
		s.useJoin = false
		return s
	}
	on := s.Joins[n].ArrayJson("on")
	for _, v := range on {
		for k := range v {
			if k == conector {
				and := v.ArrayJson(conector)
				and = append(and, cond.ToJson())
				v[k] = and
			}
		}
	}
	on = append(on, cond.ToJson())
	s.Joins[n]["on"] = on

	return s
}

/**
* And
* @param cond Condition
* @return *Cmd
**/
func (s *Ql) And(cond Condition) *Ql {
	if !s.useJoin {
		s.where.And(cond)
		return s
	}

	s.whereJoin(cond, "and")
	return s
}

/**
* Or
* @param cond Condition
* @return *Cmd
**/
func (s *Ql) Or(cond Condition) *Ql {
	if !s.useJoin {
		s.where.Or(cond)
		return s
	}

	s.whereJoin(cond, "or")
	return s
}

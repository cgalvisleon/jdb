package jdb

/**
* GroupBy
* @param fields ...string
* @return *Ql
**/
func (s *Ql) GroupBy(fields ...string) *Ql {
	for _, field := range fields {
		field := s.getField(field, false)
		if field != nil {
			s.Groups = append(s.Groups, field)
		}
	}

	return s
}

/**
* setGroupBy
* @param fields ...string
* @return *Ql
**/
func (s *Ql) setGroupBy(fields ...string) *Ql {
	return s.GroupBy(fields...)
}

/**
* listGroups
* @return []string
**/
func (s *Ql) listGroups() []string {
	result := []string{}
	for _, field := range s.Groups {
		def := s.asField(field)
		result = append(result, def)
	}

	return result
}

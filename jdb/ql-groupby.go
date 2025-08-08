package jdb

/**
* GroupBy
* @param fields ...string
* @return *Ql
**/
func (s *Ql) GroupBy(fields ...string) *Ql {
	for _, field := range fields {
		field := s.getField(field)
		if field != nil {
			s.Groups = append(s.Groups, field)
		}
	}

	return s
}

/**
* SetGroupBy
* @param fields ...string
* @return *Ql
**/
func (s *Ql) SetGroupBy(fields ...string) *Ql {
	if len(fields) == 0 {
		return s
	}

	return s.GroupBy(fields...)
}

/**
* getGroupsBy
* @return []string
**/
func (s *Ql) getGroupsBy() []string {
	result := []string{}
	for _, field := range s.Groups {
		def := s.asField(field)
		result = append(result, def)
	}

	return result
}

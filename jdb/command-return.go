package jdb

/**
* Return
* @param fields ...string
* @return *Command
**/
func (s *Command) Return(fields ...string) *Command {
	for _, name := range fields {
		field := s.getField(name)
		if field == nil {
			continue
		}

		s.Returns = append(s.Returns, field)
	}

	return s
}

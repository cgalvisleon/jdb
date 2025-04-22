package jdb

func (s *Command) Return(fields ...string) *Command {
	for _, name := range fields {
		field := s.getField(name)
		s.Returns = append(s.Returns, field)
	}

	return s
}

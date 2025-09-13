package jdb

/**
* Return
* @param fields ...string
* @return *Command
**/
func (s *Command) Return(fields ...string) *Command {
	model := s.From
	for _, pK := range model.PrimaryKeys {
		field := pK.GetField()
		s.Returns = append(s.Returns, field)
	}

	for _, name := range fields {
		field := s.getField(name)
		if field == nil {
			continue
		}

		s.Returns = append(s.Returns, field)
	}

	return s
}

/**
* SetReturn
* @param fields []string
* @return *Command
**/
func (s *Command) SetReturn(fields []string) *Command {
	return s.Return(fields...)
}

/**
* Returning
* @param fields []string
* @return *Command
**/
func (s *Command) Returning(fields []string) *Command {
	return s.Return(fields...)
}

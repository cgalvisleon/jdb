package jdb

type FullText []*Field

func (s *Model) FullText(fields ...string) FullText {
	result := FullText{}
	for _, name := range fields {
		field := s.GetField(name, false)
		if field != nil {
			result = append(result, field)
		}
	}

	return result
}

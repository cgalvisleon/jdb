package jdb

func (s *Command) beforeInsert(item map[string]*Field) map[string]*Field {
	if s.From == nil {
		return item
	}

	if s.From.IndexField != nil {
		index := s.From.GetSerie()
		field := s.From.IndexField.GetField()
		if field != nil {
			field.Value = index
			item[field.Name] = field
		}
	}

	return item
}

func (s *Command) prepare() {
	from := s.From
	for _, data := range s.Data {
		item := make(map[string]*Field, 0)
		for k, v := range data {
			field := from.GetField(k)
			if field == nil {
				continue
			}
			field.Value = v
			item[field.Name] = field
		}
		if s.Command == Insert {
			item = s.beforeInsert(item)
		}

		s.Values = append(s.Values, item)
	}

	return
}

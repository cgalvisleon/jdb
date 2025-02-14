package jdb

func setValue(value *Value, col *Column, v interface{}) *Value {
	switch col.TypeColumn {
	case TpAtribute:
		if value.Atribs[col.Name] == nil {
			value.Atribs[col.Name] = v
			value.Data[col.Name] = v
		}
	case TpColumn:
		if value.Columns[col.Name] == nil {
			value.Columns[col.Name] = v
			value.Data[col.Name] = v
		}
	}

	return value
}

func (s *Command) beforeInsert(value *Value) *Value {
	if s.From == nil {
		return value
	}

	if s.From.IndexField != nil {
		index := s.From.GetSerie()
		setValue(value, s.From.IndexField, index)
	}

	return value
}

func (s *Command) prepare() []*Value {
	from := s.From
	for _, data := range s.Origin {
		value := NewValue()
		for k, v := range data {
			field := from.GetField(k)
			if field == nil {
				if from.SourceField != nil && !from.Integrity {
					value.Atribs[k] = v
					value.Data[k] = v
				}
			} else if field.Column == from.FullTextField {
				continue
			} else {
				setValue(value, field.Column, v)
			}
		}
		if s.Command == Insert {
			value = s.beforeInsert(value)
		}
		s.Values = append(s.Values, value)
	}

	return s.Values
}

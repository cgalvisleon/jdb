package jdb

import (
	"github.com/cgalvisleon/et/utility"
)

func (s *Command) beforeInsert(item map[string]*Field) map[string]*Field {
	if s.From == nil {
		return item
	}

	if s.From.UseCore && s.From.IndexField != nil {
		field := s.From.IndexField.GetField()
		if field != nil {
			field.Value = utility.GenIndex()
			item[field.Name] = field
		}
	}

	if s.From.UseCore && s.From.SystemKeyField != nil {
		field := s.From.SystemKeyField.GetField()
		if field != nil {
			field.Value = s.From.GenId()
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
			field := from.getField(k, true)
			if field == nil {
				continue
			}
			field.setValue(v)
			item[field.Name] = field
		}
		if s.Command == Insert {
			item = s.beforeInsert(item)
		}

		s.Values = append(s.Values, item)
	}
}

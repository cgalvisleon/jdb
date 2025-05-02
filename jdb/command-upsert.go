package jdb

import (
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) upsert() error {
	if err := s.prepare(); err != nil {
		return err
	}
	model := s.From

	ql := From(model)
	for _, value := range s.Values {
		for k := range model.Required {
			field := value[k]
			if field == nil {
				return mistake.Newf(MSG_FIELD_REQUIRED_RELATION, k, model.Name)
			}

			s.Where(k).Eq(field.Value)
			ql.Where(k).Eq(field.Value)
		}
	}

	exist, err := ql.
		setDebug(s.IsDebug).
		ItExistsTx(s.tx)
	if err != nil {
		return err
	}

	if exist {
		s.Command = Update
		return s.updated()
	}

	s.Command = Insert
	return s.inserted()
}

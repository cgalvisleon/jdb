package jdb

import (
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) upsert() error {
	s.prepare()
	model := s.From

	ql := From(model)
	for _, value := range s.Values {
		for k := range model.Required {
			field := value[k]
			if field == nil {
				return mistake.Newf(MSG_FIELD_REQUIRED, k)
			}

			ql.Where(k).Eq(field.Value)
		}
	}

	exist, err := ql.ExistTx(s.tx)
	if err != nil {
		return err
	}

	if exist {
		return s.updated()
	}

	return s.inserted()
}

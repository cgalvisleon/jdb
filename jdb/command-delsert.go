package jdb

import "github.com/cgalvisleon/et/mistake"

func (s *Command) delsert() error {
	if err := s.prepare(); err != nil {
		return err
	}

	if len(s.Data) != 1 {
		return mistake.New(MSG_MANY_INSERT_DATA)
	}

	model := s.From

	ql := From(model)
	data := s.Data[0]
	where, err := model.GetWhereByRequired(data)
	if err != nil {
		return err
	}

	s.setWheres(where)
	ql.setWheres(where)

	exist, err := ql.
		setDebug(s.IsDebug).
		ItExistsTx(s.tx)
	if err != nil {
		return err
	}

	if exist {
		s.Command = Update
		return s.deleted()
	}

	s.Command = Insert
	return s.inserted()
}

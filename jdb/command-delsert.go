package jdb

import "github.com/cgalvisleon/et/mistake"

func (s *Command) delsert() error {
	if len(s.Data) != 1 {
		return mistake.New(MSG_MANY_INSERT_DATA)
	}

	model := s.From
	data := s.Data[0]
	where, err := model.GetWhereByPrimaryKeys(data)
	if err != nil {
		return err
	}

	s.current(where)
	if s.Current.Ok {
		s.SetWheres(where)
		s.Command = Delete
		return s.deleted()
	}

	s.Command = Insert
	return s.inserted()
}

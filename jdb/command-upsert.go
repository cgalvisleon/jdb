package jdb

import "fmt"

func (s *Command) upsert() error {
	if len(s.Data) != 1 {
		return fmt.Errorf(MSG_MANY_INSERT_DATA)
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
		s.Command = Update
		return s.updated()
	}

	s.Command = Insert
	return s.inserted()
}

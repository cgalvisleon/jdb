package jdb

import "fmt"

func (s *Command) upsert() error {
	model := s.getModel()
	if model == nil {
		return fmt.Errorf(MSG_MODEL_REQUIRED)
	}

	if len(s.Data) != 1 {
		return fmt.Errorf(MSG_MANY_INSERT_DATA)
	}

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

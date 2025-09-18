package jdb

import (
	"fmt"
)

func (s *Command) inserted() error {
	model := s.getModel()
	if model == nil {
		return fmt.Errorf(MSG_MODEL_REQUIRED)
	}

	if len(s.Data) == 0 {
		return fmt.Errorf(MSG_NOT_DATA, s.Command.Str(), model.Name)
	}

	if err := s.prepare(); err != nil {
		return err
	}

	results, err := s.Db.Command(s)
	if err != nil {
		return err
	}

	s.Result = results
	if !s.Result.Ok {
		return fmt.Errorf(MSG_NOT_INSERT_DATA)
	}

	s.ResultMap, err = model.getMapResultByPk(s.Result.Result)
	if err != nil {
		return err
	}

	for _, after := range s.ResultMap {
		for _, fn := range model.afterInsert {
			err := fn(s.tx, after)
			if err != nil {
				return err
			}
		}
		for _, fn := range s.afterInsert {
			err := fn(s.tx, after)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

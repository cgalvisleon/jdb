package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) updated() error {
	if len(s.Data) == 0 {
		return mistake.Newf(MSG_NOT_DATA, s.Command.Str(), s.From.Name)
	}

	if err := s.prepare(); err != nil {
		return err
	}

	model := s.From
	results, err := s.Db.Command(s)
	if err != nil {
		return err
	}

	s.Result = results
	if !s.Result.Ok {
		return mistake.New(MSG_NOT_UPDATE_DATA)
	}

	s.ResultMap, err = model.getMapResultByPk(s.Result.Result)
	if err != nil {
		return err
	}

	for key, after := range s.ResultMap {
		before := s.CurrentMap[key]
		if before == nil {
			before = et.Json{}
		}
		for _, fn := range model.afterUpdate {
			err := fn(s.tx, after)
			if err != nil {
				return err
			}
		}
		for _, fn := range s.afterUpdate {
			err := fn(s.tx, after)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

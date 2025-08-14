package jdb

import (
	"github.com/cgalvisleon/et/et"
)

func (s *Command) deleted() error {
	if err := s.prepare(); err != nil {
		return err
	}

	model := s.From
	results, err := s.Db.Command(s)
	if err != nil {
		publishError(model, s.Sql, err)
		return err
	}

	s.Result = results
	s.ResultMap, err = model.getMapResultByPk(s.Result.Result)
	if err != nil {
		return err
	}

	if !s.isSync && model.UseCore {
		publishDelete(model, s.Sql)
	}

	for _, before := range s.ResultMap {
		for _, event := range model.eventsDelete {
			err := event(s.tx, model, before, et.Json{})
			if err != nil {
				return err
			}
		}

		for _, fn := range s.afterDelete {
			err := fn(s.tx, before)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

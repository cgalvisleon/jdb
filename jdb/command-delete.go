package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
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
	if !results.Ok {
		return mistake.New(MSG_NOT_DELETE_DATA)
	}

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

		for _, jsCode := range s.afterFuncDelete {
			s.vm.Set("tx", s.tx)
			s.vm.Set("data", before)
			_, err := s.vm.RunString(jsCode)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

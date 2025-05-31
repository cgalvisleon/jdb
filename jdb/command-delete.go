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
		for _, event := range model.eventError {
			event(model, et.Json{
				"command": "delete",
				"sql":     s.Sql,
				"where":   s.getWheres(),
				"error":   err.Error(),
			})
		}

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
		audit("delete", s.Sql)
		model.Emit(EVENT_MODEL_SYNC, et.Json{
			"command": "delete",
			"db":      model.Db.Name,
			"schema":  model.Schema,
			"model":   model.Name,
			"sql":     s.Sql,
			"where":   s.getWheres(),
		})
	}

	if !model.isAudit {
		audit("delete", s.Sql)
	}

	for _, before := range s.ResultMap {
		for _, event := range model.eventsDelete {
			err := event(s.tx, model, before, et.Json{})
			if err != nil {
				return err
			}
		}

		for _, jsCode := range model.EventsDelete {
			model.vm.Set("tx", s.tx)
			model.vm.Set("before", before)
			model.vm.Set("after", et.Json{})
			_, err := model.vm.RunString(jsCode)
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

		for _, jsCode := range s.afterVmDelete {
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

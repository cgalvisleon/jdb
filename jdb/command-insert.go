package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) inserted() error {
	if len(s.Data) == 0 {
		return mistake.Newf(MSG_NOT_DATA, s.Command.Str(), s.From.Name)
	}

	if err := s.prepare(); err != nil {
		return err
	}

	model := s.From
	results, err := s.Db.Command(s)
	if err != nil {
		for _, event := range model.eventError {
			event(model, et.Json{
				"command": "insert",
				"sql":     s.Sql,
				"data":    s.Data,
				"error":   err.Error(),
			})
		}

		return err
	}

	s.Result = results
	if !s.Result.Ok {
		return mistake.New(MSG_NOT_INSERT_DATA)
	}

	s.ResultMap, err = model.getMapResultByPk(s.Result.Result)
	if err != nil {
		return err
	}

	if !s.isSync && model.UseCore {
		model.Emit(EVENT_MODEL_SYNC, et.Json{
			"command": "insert",
			"db":      model.Db.Name,
			"schema":  model.Schema,
			"model":   model.Name,
			"sql":     s.Sql,
			"values":  s.Values,
		})
	}

	for _, after := range s.ResultMap {
		for _, event := range model.eventsInsert {
			err := event(s.tx, model, et.Json{}, after)
			if err != nil {
				return err
			}
		}

		for _, jsCode := range model.EventsInsert {
			model.vm.Set("tx", s.tx)
			model.vm.Set("before", et.Json{})
			model.vm.Set("after", after)
			_, err := model.vm.RunString(jsCode)
			if err != nil {
				return err
			}
		}
	}

	for _, data := range s.Data {
		for _, fn := range s.afterInsert {
			err := fn(s.tx, data)
			if err != nil {
				return err
			}
		}

		for _, jsCode := range s.afterVmInsert {
			s.vm.Set("tx", s.tx)
			s.vm.Set("data", data)
			_, err := s.vm.RunString(jsCode)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

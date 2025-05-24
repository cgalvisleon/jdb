package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
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
		syncChannel := strs.Format("sync:%s", model.Db.Name)
		event.Publish(syncChannel, et.Json{
			"fromId":  model.Db.Id,
			"command": "delete",
			"model":   model.Name,
			"sql":     s.Sql,
			"where":   s.getWheres(),
			"result":  s.Result,
		})
	}

	for _, before := range s.ResultMap {
		for _, event := range model.eventsDelete {
			err := event(s.tx, model, before, et.Json{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

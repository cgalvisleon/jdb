package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
)

func (s *Command) deleted() error {
	model := s.From
	results, err := s.Db.Command(s)
	if err != nil {
		for _, event := range model.eventError {
			event(model, et.Json{
				"command": "delete",
				"sql":     s.Sql,
				"where":   s.listWheres(),
				"error":   err.Error(),
			})
		}

		return err
	}

	s.Result = results
	if !results.Ok {
		return mistake.New(MSG_NOT_DELETE_DATA)
	}

	if !s.isSync && model.UseCore {
		syncChannel := strs.Format("sync:%s", model.Db.Name)
		event.Publish(syncChannel, et.Json{
			"fromId":  model.Db.Id,
			"command": "delete",
			"model":   model.Name,
			"sql":     s.Sql,
			"where":   s.listWheres(),
			"result":  s.Result,
		})
	}

	if s.isUndo {
		return nil
	}

	for _, result := range s.Result.Result {
		before := result.ValJson(et.Json{}, "result", "before")
		after := result.ValJson(et.Json{}, "result", "after")

		for _, event := range model.eventsDelete {
			err := event(s.tx, model, before, after)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

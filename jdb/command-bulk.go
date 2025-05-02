package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/strs"
)

func (s *Command) bulk() error {
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
	if !results.Ok {
		return nil
	}

	if model.UseCore && model.SystemKeyField != nil {
		syncChannel := strs.Format("sync:%s", model.Db.Name)
		event.Publish(syncChannel, et.Json{
			"fromId":  model.Db.Id,
			"command": "insert",
			"sql":     s.Sql,
			"values":  s.Values,
			"result":  s.Result,
		})
	}

	if s.isUndo {
		return nil
	}

	for _, result := range s.Result.Result {
		before := result.ValJson(et.Json{}, "result", "before")
		after := result.ValJson(et.Json{}, "result", "after")

		for _, event := range model.eventsInsert {
			err := event(s.tx, model, before, after)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

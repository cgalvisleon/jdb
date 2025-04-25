package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/strs"
)

func (s *Command) inserted() error {
	s.prepare()
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

	if s.rollback {
		return nil
	}

	for _, result := range s.Result.Result {
		before := result.ValJson(et.Json{}, "result", "before")
		after := result.ValJson(et.Json{}, "result", "after")

		go func() {
			if model.SystemKeyField != nil {
				sysid := after.Str(model.SystemKeyField.Name)
				s.Db.upsertRecord(model.Table, "insert", sysid)
			}
		}()

		for _, event := range model.eventsInsert {
			err := event(model, before, after)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

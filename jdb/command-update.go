package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/strs"
)

func (s *Command) updated() error {
	s.prepare()
	model := s.From

	results, err := s.Db.Command(s)
	if err != nil {
		for _, event := range model.eventError {
			event(model, et.Json{
				"command": "update",
				"sql":     s.Sql,
				"data":    s.Data,
				"where":   s.listWheres(),
				"error":   err.Error(),
			})
		}

		return err
	}

	s.Result = results
	if !s.Result.Ok {
		return nil
	}

	if model.UseCore {
		syncChannel := strs.Format("sync:%s", model.Db.Name)
		event.Publish(syncChannel, et.Json{
			"fromId":  model.Db.Id,
			"command": "update",
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
				s.Db.upsertRecord(model.Table, "update", sysid)
			}
		}()

		for _, event := range model.eventsUpdate {
			err := event(model, before, after)
			if err != nil {
				return err
			}
		}

		changed := before.IsChanged(after)
		if s.history && changed && model.History != nil {
			err := eventHistoryDefault(model, before)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

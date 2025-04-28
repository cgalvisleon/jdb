package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/strs"
)

func (s *Command) updated(tx *sql.Tx) error {
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

	if s.isUndo {
		return nil
	}

	for _, result := range s.Result.Result {
		before := result.ValJson(et.Json{}, "result", "before")
		after := result.ValJson(et.Json{}, "result", "after")
		changed := before.IsChanged(after)

		if !changed {
			continue
		}

		for _, event := range model.eventsUpdate {
			err := event(tx, model, before, after)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

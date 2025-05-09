package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
)

func (s *Command) updated() error {
	if err := s.prepare(); err != nil {
		return err
	}
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
		return mistake.New(MSG_NOT_UPDATE_DATA)
	}

	if !s.isSync && model.UseCore {
		syncChannel := strs.Format("sync:%s", model.Db.Name)
		event.Publish(syncChannel, et.Json{
			"fromId":  model.Db.Id,
			"command": "update",
			"model":   model.Name,
			"sql":     s.Sql,
			"values":  s.Values,
			"result":  s.Result,
		})
	}

	for _, result := range s.Result.Result {
		before := result.ValJson(et.Json{}, "result", "before")
		after := result.ValJson(et.Json{}, "result", "after")
		changed := before.IsChanged(after)

		if !changed {
			continue
		}

		for _, event := range model.eventsUpdate {
			err := event(s.tx, model, before, after)
			if err != nil {
				return err
			}
		}
	}

	for _, data := range s.Data {
		for _, fn := range s.afterUpdate {
			err := fn(s.tx, data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

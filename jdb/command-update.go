package jdb

import (
	"github.com/cgalvisleon/et/et"
)

func (s *Command) updated() error {
	s.prepare()
	model := s.From

	results, err := s.Db.Command(s)
	if err != nil {
		for _, event := range s.From.EventError {
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

	if s.rollback {
		return nil
	}

	for _, result := range s.Result.Result {
		before := result.ValJson(et.Json{}, "result", "before")
		after := result.ValJson(et.Json{}, "result", "after")

		for _, event := range s.From.EventsUpdate {
			err := event(model, before, after)
			if err != nil {
				return err
			}
		}

		changed := before.IsChanged(after)
		if s.history && changed && model.History.With != nil {
			err := EventHistoryDefault(model, before, after)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

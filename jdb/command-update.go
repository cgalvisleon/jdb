package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) update() error {
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

	if !results.Ok {
		return mistake.New(MSG_NOT_UPDATE_DATA)
	}

	s.Result = results

	return nil
}

func (s *Command) updated() error {
	err := s.update()
	if err != nil {
		return err
	}

	model := s.From
	for _, result := range s.Result.Result {
		before := result.ValJson(et.Json{}, "result", "before")
		after := result.ValJson(et.Json{}, "result", "after")

		for _, event := range s.From.EventsUpdate {
			err := event(model, before, after)
			if err != nil {
				return err
			}
		}

		if s.history && model.History.With != nil {
			err := EventHistoryDefault(model, before, after)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

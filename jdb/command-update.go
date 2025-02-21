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
		if model.EventError != nil {
			model.EventError(model, et.Json{
				"command": "update",
				"sql":     s.Sql,
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
		before := result.Json("before")
		after := result.Json("after")

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

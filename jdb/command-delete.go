package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) delete() error {
	model := s.From
	results, err := s.Db.Command(s)
	if err != nil {
		if model.EventError != nil {
			model.EventError(model, et.Json{
				"command": "insert",
				"sql":     s.Sql,
				"error":   err.Error(),
			})
		}
		return err
	}

	if !results.Ok {
		return mistake.New(MSG_NOT_DELETE_DATA)
	}

	s.Result = results

	for _, result := range results.Result {
		before := result.Json("before")
		after := result.Json("after")

		for _, event := range s.From.EventsUpdate {
			err := event(model, before, after)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

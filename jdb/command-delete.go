package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) delete() error {
	model := s.From
	results, err := s.Db.Command(s)
	if err != nil {
		for _, event := range s.From.EventError {
			event(model, et.Json{
				"command": "insert",
				"sql":     s.Sql,
				"where":   s.listWheres(),
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
		before := result.ValJson(et.Json{}, "result", "before")
		after := result.ValJson(et.Json{}, "result", "after")

		for _, event := range s.From.EventsUpdate {
			err := event(model, before, after)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

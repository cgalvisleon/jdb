package jdb

import (
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) delete() error {
	results, err := s.Db.Command(s)
	if err != nil {
		return err
	}

	if !results.Ok {
		return mistake.New(MSG_NOT_DELETE_DATA)
	}

	s.Result = results

	model := s.From.Model
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

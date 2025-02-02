package jdb

import (
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) bulk() error {
	s.prepare()

	results, err := s.Db.Command(s)
	if err != nil {
		return err
	}

	if !results.Ok {
		return mistake.New(MSG_NOT_INSERT_DATA)
	}

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

	s.Result = results

	return nil
}

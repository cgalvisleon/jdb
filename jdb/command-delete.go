package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) delete() (et.Items, error) {
	results, err := s.Db.Command(s)
	if err != nil {
		return et.Items{}, err
	}

	if !results.Ok {
		return et.Items{}, mistake.New(MSG_NOT_DELETE_DATA)
	}

	model := s.From.Model
	for _, result := range results.Result {
		before := result.Json("before")
		after := result.Json("after")

		for _, event := range s.From.EventsUpdate {
			err := event(model, before, after)
			if err != nil {
				return et.Items{}, err
			}
		}
	}

	return results, nil
}

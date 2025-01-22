package jdb

import (
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) update() error {
	s.consolidate()

	results, err := s.Db.Command(s)
	if err != nil {
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

	model := s.From.Model
	for _, result := range s.Result.Result {
		before := result.Json("before")
		after := result.Json("after")

		for _, event := range s.From.EventsUpdate {
			err := event(model, before, after)
			if err != nil {
				return err
			}
		}

		s.Commit = append(s.Commit, after)

		if model.History != nil {
			err := EventHistoryDefault(model, before, after)
			if err != nil {
				return err
			}
		}
	}

	for _, command := range s.Commands {
		_, err := command.Exec()
		if err != nil {
			break
		}
	}

	return nil
}

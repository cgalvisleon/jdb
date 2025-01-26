package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/utility"
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

	changed := func(before, after et.Json, exclude []string) bool {
		for k, v := range before {
			if slices.Contains(exclude, k) {
				continue
			}

			if utility.Quote(after[k]) != utility.Quote(v) {
				return false
			}
		}

		return true
	}

	model := s.From.Model
	for _, result := range s.Result.Result {
		before := result.Json("before")
		after := result.Json("after")

		if !changed(before, after, []string{CREATED_AT, UPDATED_AT}) {
			continue
		}

		for _, event := range s.From.EventsUpdate {
			err := event(model, before, after)
			if err != nil {
				return err
			}
		}

		if model.History != nil {
			err := EventHistoryDefault(model, before, after)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

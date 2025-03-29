package jdb

import "github.com/cgalvisleon/et/et"

func (s *Command) bulk() error {
	s.prepare()
	model := s.From

	results, err := s.Db.Command(s)
	if err != nil {
		for _, event := range s.From.EventError {
			event(model, et.Json{
				"command": "insert",
				"sql":     s.Sql,
				"data":    s.Data,
				"error":   err.Error(),
			})
		}

		return err
	}

	s.Result = results
	if !results.Ok {
		return nil
	}

	if s.rollback {
		return nil
	}

	for _, result := range s.Result.Result {
		before := result.ValJson(et.Json{}, "result", "before")
		after := result.ValJson(et.Json{}, "result", "after")

		for _, event := range s.From.EventsInsert {
			err := event(model, before, after)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Command) inserted() error {
	err := s.bulk()
	if err != nil {
		return err
	}

	return nil
}

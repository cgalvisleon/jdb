package jdb

import "github.com/cgalvisleon/et/et"

func (s *Command) updated(old, data et.Json) (et.Item, error) {
	s.consolidate(data)

	for _, trigger := range s.Model.BeforeUpdate {
		err := Triggers[trigger](old, s.New, data)
		if err != nil {
			return et.Item{}, err
		}
	}

	result, err := s.Db.Command(s)
	if err != nil {
		return et.Item{}, err
	}

	if result.Ok {
		s.New = &result.Result
	}

	for _, trigger := range s.Model.AfterUpdate {
		err := Triggers[trigger](old, s.New, data)
		if err != nil {
			return et.Item{}, err
		}
	}

	s.Model.GetDetails(&result.Result)

	return result, nil
}

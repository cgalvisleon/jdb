package jdb

import "github.com/cgalvisleon/et/et"

func (s *Command) updated(old, data et.Json) (et.Item, error) {
	s.consolidate(data)

	for _, trigger := range s.BeforeUpdate {
		err := Triggers[trigger](old, s.New, data)
		if err != nil {
			return et.Item{}, err
		}
	}

	result, err := s.Db.Command(s)
	if err != nil {
		return et.Item{}, err
	} else {
		result.Ok = true
	}

	if result.Ok {
		s.New = &result.Result
	}

	for _, trigger := range s.AfterUpdate {
		err := Triggers[trigger](old, s.New, data)
		if err != nil {
			return et.Item{}, err
		}
	}

	if len(s.Returns) > 0 {
		s.GetDetails(&result.Result)
	}

	return result, nil
}

package jdb

import "github.com/cgalvisleon/et/et"

func (s *Command) inserted(data et.Json) (et.Item, error) {
	s.consolidate(data)

	for _, trigger := range s.BeforeInsert {
		err := Triggers[trigger](nil, s.New, data)
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

	for _, trigger := range s.AfterInsert {
		err := Triggers[trigger](nil, s.New, data)
		if err != nil {
			return et.Item{}, err
		}
	}

	if len(s.Returns) > 0 {
		s.GetDetails(&result.Result)
	}

	return result, nil
}

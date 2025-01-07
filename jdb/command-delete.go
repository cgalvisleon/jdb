package jdb

import "github.com/cgalvisleon/et/et"

func (s *Command) delete(old et.Json) (et.Item, error) {
	for _, trigger := range s.Model.BeforeDelete {
		err := Triggers[trigger](old, nil, nil)
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

	for _, trigger := range s.Model.AfterDelete {
		err := Triggers[trigger](old, nil, nil)
		if err != nil {
			return et.Item{}, err
		}
	}

	s.Model.GetDetails(&result.Result)

	return result, nil
}

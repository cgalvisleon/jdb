package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) updated() (et.Items, error) {
	current, err := s.Db.Current(s)
	if err != nil {
		return et.Items{}, err
	}
	data := s.Origin[0]
	s.consolidate(data)
	results := et.Items{}

	if !current.Ok {
		return et.Items{}, mistake.New(MSG_NOT_DATA)
	}

	if current.Count > 1 {
		result, err := s.Db.Command(s)
		if err != nil {
			return et.Items{}, err
		}

		new := &result.Result
		s.From.GetDetails(new)
		results.Add(*new)

		return results, nil
	}

	old := current.Result[0]
	for _, trigger := range s.From.BeforeUpdate {
		err := Triggers[trigger](old, s.New, data)
		if err != nil {
			return et.Items{}, err
		}
	}

	result, err := s.Db.Command(s)
	if err != nil {
		return et.Items{}, err
	}

	new := &result.Result
	for _, trigger := range s.From.AfterUpdate {
		err := Triggers[trigger](old, new, data)
		if err != nil {
			return et.Items{}, err
		}
	}

	s.From.GetDetails(new)
	results.Add(*new)

	return results, nil
}

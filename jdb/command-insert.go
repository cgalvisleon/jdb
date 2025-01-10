package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) inserted() (et.Item, error) {
	if s.From.SystemKeyField == nil {
		return et.Item{}, mistake.New(MSG_SYSTEMKEYFIELD_NOT_FOUND)
	}

	data := s.Origin[0]
	s.consolidate(data)

	for _, trigger := range s.From.BeforeInsert {
		err := Triggers[trigger](nil, s.New, data)
		if err != nil {
			return et.Item{}, err
		}
	}

	result, err := s.Db.Command(s)
	if err != nil {
		return et.Item{}, err
	}

	new := &result.Result

	for _, trigger := range s.From.AfterInsert {
		err := Triggers[trigger](nil, new, data)
		if err != nil {
			return et.Item{}, err
		}
	}

	s.From.GetDetails(new)

	return result, nil
}

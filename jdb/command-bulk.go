package jdb

import "github.com/cgalvisleon/et/et"

func (s *Command) bulk() (et.Item, error) {
	result, err := s.Db.Command(s)
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}

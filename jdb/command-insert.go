package jdb

import (
	"github.com/cgalvisleon/et/et"
)

func (s *Command) inserted() (et.Item, error) {
	result, err := s.bulk()
	if err != nil {
		return et.Item{}, err
	}

	return et.Item{
		Ok:     true,
		Result: result.Result[0],
	}, nil
}

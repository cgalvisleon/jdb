package jdb

import (
	"github.com/cgalvisleon/et/et"
)

/**
* insert
* @return error
**/
func (s *Cmd) insert() (et.Items, error) {
	for _, new := range s.Data {
		for _, fn := range s.beforeInserts {
			err := fn(s.tx, et.Json{}, new)
			if err != nil {
				return et.Items{}, err
			}
		}

		result, err := s.db.command(s)
		if err != nil {
			return et.Items{}, err
		}

		if !result.Ok {
			continue
		}

		new = result.First().Result
		s.Result.Add(new)

		for _, fn := range s.afterInserts {
			err := fn(s.tx, et.Json{}, new)
			if err != nil {
				return et.Items{}, err
			}
		}
	}

	return s.Result, nil
}

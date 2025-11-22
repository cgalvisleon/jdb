package jdb

import (
	"github.com/cgalvisleon/et/et"
)

/**
* insert
* @return error
**/
func (s *Cmd) insert() (et.Items, error) {
	for _, data := range s.Data {
		for _, fn := range s.beforeInserts {
			err := fn(s.tx, data)
			if err != nil {
				return et.Items{}, err
			}
		}

		result, err := s.db.command(s)
		if err != nil {
			return et.Items{}, err
		}

		data := result.First().Result
		s.Result.Add(data)
		if !result.Ok {
			return result, nil
		}

		for _, fn := range s.afterInserts {
			err := fn(s.tx, data)
			if err != nil {
				return et.Items{}, err
			}
		}
	}

	return s.Result, nil
}

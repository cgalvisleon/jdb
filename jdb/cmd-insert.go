package jdb

import (
	"github.com/cgalvisleon/et/et"
)

/**
* insertTx
* @return error
**/
func (s *Cmd) insertTx() (et.Items, error) {
	for _, data := range s.Data {
		for _, definition := range s.BeforeInserts {
			s.vm.Set("data", data)
			_, err := s.vm.Run(definition)
			if err != nil {
				return et.Items{}, err
			}
		}

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

		for _, definition := range s.AfterInserts {
			s.vm.Set("data", data)
			_, err := s.vm.Run(definition)
			if err != nil {
				return et.Items{}, err
			}
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

package jdb

import "github.com/cgalvisleon/et/et"

func (s *Cmd) updateTx() (et.Items, error) {
	current, err := s.From.
		Query(et.Json{
			"where": s.Wheres,
		}).
		Debug().
		All()
	if err != nil {
		return et.Items{}, err
	}

	for _, data := range current.Result {
		for k, v := range s.Data[0] {
			data[k] = v
		}

		for _, definition := range s.BeforeUpdates {
			s.vm.Set("data", data)
			_, err := s.vm.Run(definition)
			if err != nil {
				return et.Items{}, err
			}
		}

		for _, fn := range s.beforeUpdates {
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

		for _, definition := range s.AfterUpdates {
			s.vm.Set("data", data)
			_, err := s.vm.Run(definition)
			if err != nil {
				return et.Items{}, err
			}
		}

		for _, fn := range s.afterUpdates {
			err := fn(s.tx, data)
			if err != nil {
				return et.Items{}, err
			}
		}
	}

	return s.Result, nil
}

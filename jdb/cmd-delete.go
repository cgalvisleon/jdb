package jdb

import "github.com/cgalvisleon/et/et"

func (s *Cmd) deleteTx() (et.Items, error) {
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
		for _, definition := range s.beforeDeletes {
			s.vm.Set("data", data)
			_, err := s.vm.Run(definition)
			if err != nil {
				return et.Items{}, err
			}
		}

		for _, fn := range s.eventBeforeDelete {
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

		for _, definition := range s.afterDeletes {
			s.vm.Set("data", data)
			_, err := s.vm.Run(definition)
			if err != nil {
				return et.Items{}, err
			}
		}

		for _, fn := range s.eventAfterDelete {
			err := fn(s.tx, data)
			if err != nil {
				return et.Items{}, err
			}
		}
	}

	return et.Items{}, nil
}

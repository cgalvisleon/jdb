package jdb

import "github.com/cgalvisleon/et/et"

/**
* delete
* @return (et.Items, error)
**/
func (s *Cmd) delete() (et.Items, error) {
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
		for _, definition := range s.BeforeDeletes {
			s.vm.Set("data", data)
			_, err := s.vm.Run(definition)
			if err != nil {
				return et.Items{}, err
			}
		}

		for _, fn := range s.beforeDeletes {
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

		for _, definition := range s.AfterDeletes {
			s.vm.Set("data", data)
			_, err := s.vm.Run(definition)
			if err != nil {
				return et.Items{}, err
			}
		}

		for _, fn := range s.afterDeletes {
			err := fn(s.tx, data)
			if err != nil {
				return et.Items{}, err
			}
		}
	}

	return et.Items{}, nil
}

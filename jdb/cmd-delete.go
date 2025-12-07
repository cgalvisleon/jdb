package jdb

import "github.com/cgalvisleon/et/et"

/**
* delete
* @return (et.Items, error)
**/
func (s *Cmd) delete() (et.Items, error) {
	if len(s.Froms) == 0 {
		return et.Items{}, nil
	}

	from := s.Froms[0].Model
	current, err := from.
		Query(et.Json{
			"where": s.Wheres,
		}).
		All()
	if err != nil {
		return et.Items{}, err
	}

	for _, old := range current.Result {
		for _, fn := range s.beforeDeletes {
			err := fn(s.tx, old, et.Json{})
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

		old = result.First().Result
		s.Result.Add(old)

		for _, fn := range s.afterDeletes {
			err := fn(s.tx, old, et.Json{})
			if err != nil {
				return et.Items{}, err
			}
		}
	}

	return et.Items{}, nil
}

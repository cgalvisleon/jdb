package jdb

import "github.com/cgalvisleon/et/et"

/**
* update
* @return (et.Items, error)
**/
func (s *Cmd) update() (et.Items, error) {
	current, err := s.From.
		Query(et.Json{
			"where": s.Wheres,
		}).
		All()
	if err != nil {
		return et.Items{}, err
	}

	for _, old := range current.Result {
		new := old.Clone()
		for k, v := range s.Data[0] {
			new[k] = v
		}

		for _, fn := range s.beforeUpdates {
			err := fn(s.tx, old, new)
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

		for _, fn := range s.afterUpdates {
			err := fn(s.tx, old, new)
			if err != nil {
				return et.Items{}, err
			}
		}
	}

	return s.Result, nil
}

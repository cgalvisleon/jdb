package jdb

import "github.com/cgalvisleon/et/et"

/**
* update
* @return (et.Items, error)
**/
func (s *Cmd) update() (et.Items, error) {
	if len(s.Froms) == 0 {
		return et.Items{}, nil
	}

	from := s.Froms[0].Model
	current, err := from.
		WhereByCMD(s).
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

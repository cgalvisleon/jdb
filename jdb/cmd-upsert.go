package jdb

import "github.com/cgalvisleon/et/et"

/**
* upsert
* @return et.Items, error
**/
func (s *Cmd) upsert() (et.Items, error) {
	if len(s.Froms) == 0 {
		return et.Items{}, nil
	}

	from := s.Froms["from"]
	data := s.Data[0]
	keys := s.getKeys(data, "A")
	exists, err := from.
		Query(et.Json{
			"where": keys,
		}).
		Debug().
		ItExists()
	if err != nil {
		return et.Items{}, err
	}

	if exists {
		s.Command = CmdUpdate
		return s.update()
	}

	s.Command = CmdInsert
	return s.insert()
}

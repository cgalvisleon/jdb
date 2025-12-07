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

	from := s.Froms[0].Model
	data := s.Data[0]
	exists, err := from.
		WhereByPrimaryKeys(data).
		SetDebug(s.IsDebug).
		Debug().
		ItExists()
	if err != nil {
		return et.Items{}, err
	}

	if exists {
		s.Command = CmdUpdate
		s.WhereByPrimaryKeys(data)
		return s.update()
	}

	s.Command = CmdInsert
	return s.insert()
}

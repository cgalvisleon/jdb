package jdb

import "github.com/cgalvisleon/et/et"

/**
* upsertTx
* @return et.Items, error
**/
func (s *Cmd) upsertTx() (et.Items, error) {
	data := s.Data[0]
	keys := s.getKeys(data, "A")
	exists, err := s.From.
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
		return s.updateTx()
	}

	s.Command = CmdInsert
	return s.insertTx()
}

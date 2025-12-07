package jdb

import "github.com/cgalvisleon/et/et"

type Detail struct {
	From    *Model   `json:"from"`
	Fks     et.Json  `json:"fks"`
	Selects []string `json:"selects"`
}

/**
* @param tx *Tx, data et.Json
* @return et.Items, error
**/
func (s *Detail) QueryTx(tx *Tx, data et.Json) (et.Items, error) {
	ql := s.From.Select(s.Selects...)
	for fk := range s.Fks {
		kf := s.Fks.Str(fk)
		ql.Where(Eq(fk, data[kf]))
	}

	return ql.AllTx(tx)
}

/**
* @params data et.Json
* @return et.Items, error
**/
func (s *Detail) Query(data et.Json) (et.Items, error) {
	return s.QueryTx(nil, data)
}

package jdb

import "github.com/cgalvisleon/et/et"

/**
* ExecTx
* @param tx *Tx
* @return (et.Items, error)
**/
func (s *Cmd) ExecTx(tx *Tx) (et.Items, error) {
	s.tx = tx

	if err := s.validate(); err != nil {
		return et.Items{}, err
	}

	return et.Items{}, nil
}

/**
* Exec
* @return (et.Items, error)
**/
func (s *Cmd) Exec() (et.Items, error) {
	return s.ExecTx(nil)
}

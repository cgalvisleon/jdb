package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
)

/**
* queryTx
* @param tx *Tx
* @return (et.Items, error)
**/
func (s *Cmd) queryTx(tx *Tx) (et.Items, error) {
	if s.db == nil {
		return et.Items{}, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	s.setTx(tx)
	if err := s.before(); err != nil {
		return et.Items{}, err
	}

	result, err := s.db.command(s)
	if err != nil {
		return et.Items{}, err
	}

	if err := s.after(); err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* ExecTx
* @param tx *Tx
* @return (et.Items, error)
**/
func (s *Cmd) ExecTx(tx *Tx) (et.Items, error) {
	result, err := s.queryTx(tx)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* Exec
* @return (et.Items, error)
**/
func (s *Cmd) Exec() (et.Items, error) {
	return s.ExecTx(nil)
}

/**
* One
* @return (et.Item, error)
**/
func (s *Cmd) OneTx(tx *Tx) (et.Item, error) {
	result, err := s.ExecTx(tx)
	if err != nil {
		return et.Item{}, err
	}

	return result.First(), nil
}

/**
* One
* @return (et.Item, error)
**/
func (s *Cmd) One() (et.Item, error) {
	return s.OneTx(nil)
}

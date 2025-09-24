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

	err := s.validate()
	if err != nil {
		return et.Items{}, err
	}

	result, err := s.db.command(s)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* One
* @return (et.Item, error)
**/
func (s *Cmd) OneTx(tx *Tx) (et.Item, error) {
	items, err := s.queryTx(tx)
	if err != nil {
		return et.Item{}, err
	}

	return items.First(), nil
}

/**
* One
* @return (et.Item, error)
**/
func (s *Cmd) One() (et.Item, error) {
	return s.OneTx(nil)
}

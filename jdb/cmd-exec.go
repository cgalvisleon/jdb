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
	s.vm.Set("model", s.From)
	switch s.Command {
	case CmdInsert:
		return s.insertTx()
	case CmdUpdate:
		return s.updateTx()
	case CmdDelete:
		return s.deleteTx()
	case CmdUpsert:
		return s.upsertTx()
	default:
		return et.Items{}, fmt.Errorf("invalid command: %s", s.Command)
	}
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

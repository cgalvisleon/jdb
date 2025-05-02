package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

/**
* Tx
* @return *Tx
**/
func (s *Command) Tx() *Tx {
	return s.tx
}

/**
* ExecTx
* @param tx *Tx
* @return et.Items, error
**/
func (s *Command) ExecTx(tx *Tx) (et.Items, error) {
	var err error
	if tx == nil {
		tx = NewTx()

		defer func() {
			if err == nil {
				tx.Commit()
			}
		}()
	}

	s.setTx(tx)

	switch s.Command {
	case Insert:
		if len(s.Data) == 0 {
			return et.Items{}, mistake.Newf(MSG_NOT_DATA, s.Command.Str(), s.From.Name)
		}

		err := s.inserted()
		if err != nil {
			return et.Items{}, err
		}
	case Update:
		if len(s.Data) == 0 {
			return et.Items{}, mistake.Newf(MSG_NOT_DATA, s.Command.Str(), s.From.Name)
		}

		err := s.updated()
		if err != nil {
			return et.Items{}, err
		}
	case Upsert:
		err := s.upsert()
		if err != nil {
			return et.Items{}, err
		}
	case Delete:
		err := s.deleted()
		if err != nil {
			return et.Items{}, err
		}
	case Bulk:
		if len(s.Data) == 0 {
			return et.Items{}, mistake.Newf(MSG_NOT_DATA, s.Command.Str(), s.From.Name)
		}

		err := s.bulk()
		if err != nil {
			return et.Items{}, err
		}
	case Sync:
		if len(s.Data) == 0 {
			return et.Items{}, mistake.Newf(MSG_NOT_DATA, s.Command.Str(), s.From.Name)
		}

		err := s.sync()
		if err != nil {
			return et.Items{}, err
		}
	default:
		return et.Items{}, mistake.New(MSG_NOT_COMMAND)
	}

	return s.Result, nil
}

/**
* OneTx
* @param tx *Tx
* @return et.Item, error
**/
func (s *Command) OneTx(tx *Tx) (et.Item, error) {
	result, err := s.ExecTx(tx)
	if err != nil {
		return et.Item{}, err
	}

	return result.First(), nil
}

/**
* Exec
* @return et.Items, error
**/
func (s *Command) Exec() (et.Items, error) {
	return s.ExecTx(nil)
}

/**
* One
* @return et.Item, error
**/
func (s *Command) One() (et.Item, error) {
	return s.OneTx(nil)
}

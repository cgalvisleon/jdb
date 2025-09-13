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

		defer func() (et.Items, error) {
			if err == nil {
				err = tx.Commit()
				if err != nil {
					return et.Items{}, err
				}
			}

			return s.Result, err
		}()
	}

	s.setTx(tx)

	switch s.Command {
	case Insert:
		err := s.inserted()
		if err != nil {
			return et.Items{}, err
		}
	case Update:
		where := s.getWheres()
		err = s.current(where)
		if err != nil {
			return et.Items{}, err
		}
		err = s.updated()
		if err != nil {
			return et.Items{}, err
		}
	case Delete:
		err = s.deleted()
		if err != nil {
			return et.Items{}, err
		}
	case Upsert:
		err := s.upsert()
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

/**
* DoTx
* @param tx *Tx
* @return map[string]interface{}, error
**/
func (s *Command) DoTx(tx *Tx, wheres et.Json) (map[string]interface{}, error) {
	result, err := s.
		SetWhere(wheres).
		ExecTx(tx)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return result.ToMap(), nil
}

/**
* Do
* @return map[string]interface{}, error
**/
func (s *Command) Do(wheres et.Json) (map[string]interface{}, error) {
	return s.DoTx(nil, wheres)
}

/**
* DoOneTx
* @param tx *Tx
* @return map[string]interface{}, error
**/
func (s *Command) DoOneTx(tx *Tx, wheres et.Json) (map[string]interface{}, error) {
	result, err := s.
		SetWhere(wheres).
		ExecTx(tx)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return result.First().ToMap(), nil
}

/**
* DoOne
* @return map[string]interface{}, error
**/
func (s *Command) DoOne(wheres et.Json) (map[string]interface{}, error) {
	return s.DoOneTx(nil, wheres)
}

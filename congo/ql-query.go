package jdb

import "fmt"

/**
* ItExistsTx
* @param tx *Tx
* @return (bool, error)
**/
func (s *Ql) ItExistsTx(tx *Tx) (bool, error) {
	if s.db == nil {
		return false, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	s.setTx(tx)
	result, err := s.db.exists(s)
	if err != nil {
		return false, err
	}

	return result, nil
}

/**
* ItExists
* @return (bool, error)
**/
func (s *Ql) ItExists() (bool, error) {
	return s.ItExistsTx(nil)
}

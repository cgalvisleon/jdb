package jdb

import (
	"fmt"
	"sync"

	"github.com/cgalvisleon/et/et"
)

/**
* getRollupsTx
* @param tx *Tx, data et.Json
* @return
**/
func (s *Ql) getRollupsTx(tx *Tx, data et.Json) {

}

/**
* getRelationsTx
* @param tx *Tx, data et.Json
* @return
**/
func (s *Ql) getRelationsTx(tx *Tx, data et.Json) {

}

/**
* getCallsTx
* @param tx *Tx, data et.Json
* @return
**/
func (s *Ql) getCallsTx(tx *Tx, data et.Json) {

}

/**
* FirstTx
* @param tx *Tx, n int
* @return et.Items, error
**/
func (s *Ql) FirstTx(tx *Tx, n int) (et.Items, error) {
	if s.db == nil {
		return et.Items{}, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	s.setTx(tx)
	s.Limit = et.Json{
		"page": 1,
		"rows": n,
	}
	err := s.validate()
	if err != nil {
		return et.Items{}, err
	}

	result, err := s.db.query(s)
	if err != nil {
		return et.Items{}, err
	}

	wg := &sync.WaitGroup{}
	for _, data := range result.Result {
		wg.Add(1)
		go func(data et.Json) {
			defer wg.Done()
			s.getRollupsTx(tx, data)
			s.getRelationsTx(tx, data)
			s.getCallsTx(tx, data)
		}(data)
	}
	wg.Wait()

	return result, nil
}

/**
* AllTx
* @param tx *Tx
* @return et.Items, error
**/
func (s *Ql) AllTx(tx *Tx) (et.Items, error) {
	return s.FirstTx(tx, 0)
}

/**
* LastTx
* @param tx *Tx, n int
* @return et.Items, error
**/
func (s *Ql) LastTx(tx *Tx, n int) (et.Items, error) {
	return s.FirstTx(tx, n*-1)
}

/**
* OneTx
* @param tx *Tx
* @return et.Item, error
**/
func (s *Ql) OneTx(tx *Tx) (et.Item, error) {
	result, err := s.FirstTx(tx, 1)
	if err != nil {
		return et.Item{}, err
	}

	return result.First(), nil
}

/**
* RowsTx
* @param tx *Tx, limit int
* @return et.Items, error
**/
func (s *Ql) RowsTx(tx *Tx, val int) (et.Items, error) {
	return s.FirstTx(tx, val)
}

/**
* ItExistsTx
* @param tx *Tx
* @return bool, error
**/
func (s *Ql) ItExistsTx(tx *Tx) (bool, error) {
	if s.db == nil {
		return false, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	s.setTx(tx)
	err := s.validate()
	if err != nil {
		return false, err
	}
	result, err := s.db.query(s)
	if err != nil {
		return false, err
	}

	if result.Count == 0 {
		return false, nil
	}

	return result.Bool(0, "exists"), nil
}

/**
* CountedTx
* @param tx *Tx
* @return int, error
**/
func (s *Ql) CountedTx(tx *Tx) (int, error) {
	if s.db == nil {
		return 0, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	s.setTx(tx)
	err := s.validate()
	if err != nil {
		return 0, err
	}
	result, err := s.db.query(s)
	if err != nil {
		return 0, err
	}

	if result.Count == 0 {
		return 0, nil
	}

	return result.Int(0, "all"), nil
}

/**
* First
* @param n int
* @return et.Items, error
**/
func (s *Ql) First(n int) (et.Items, error) {
	return s.FirstTx(nil, n)
}

/**
* All
* @return et.Items, error
**/
func (s *Ql) All() (et.Items, error) {
	return s.AllTx(nil)
}

/**
* Last
* @param n int
* @return et.Items, error
**/
func (s *Ql) Last(n int) (et.Items, error) {
	return s.LastTx(nil, n)
}

/**
* One
* @return et.Item, error
**/
func (s *Ql) One() (et.Item, error) {
	return s.OneTx(nil)
}

/**
* Rows
* @param n int
* @return et.Items, error
**/
func (s *Ql) Rows(n int) (et.Items, error) {
	return s.RowsTx(nil, n)
}

/**
* ItExists
* @return bool, error
**/
func (s *Ql) ItExists() (bool, error) {
	return s.ItExistsTx(nil)
}

/**
* Counted
* @return int, error
**/
func (s *Ql) Counted() (int, error) {
	return s.CountedTx(nil)
}

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
func (s *Cmd) getRollupsTx(tx *Tx, data et.Json) {

}

/**
* getRelationsTx
* @param tx *Tx, data et.Json
* @return
**/
func (s *Cmd) getRelationsTx(tx *Tx, data et.Json) {

}

/**
* getCallsTx
* @param tx *Tx, data et.Json
* @return
**/
func (s *Cmd) getCallsTx(tx *Tx, data et.Json) {

}

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

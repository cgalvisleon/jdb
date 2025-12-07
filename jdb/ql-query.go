package jdb

import (
	"fmt"
	"sync"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
)

/**
* getRollupsTx
* @param tx *Tx, data et.Json
* @return
**/
func (s *Ql) getRollupsTx(tx *Tx, data et.Json) {
	for name, rollup := range s.Rollups {
		items, err := rollup.QueryTx(tx, data)
		if err != nil {
			logs.Error(err)
		}

		item := items.First().Result
		if len(item) == 1 {
			for _, v := range item {
				data[name] = v
			}
		} else {
			data[name] = item
		}
	}
}

/**
* getRelationsTx
* @param tx *Tx, data et.Json
* @return
**/
func (s *Ql) getRelationsTx(tx *Tx, data et.Json) {
	for name, relation := range s.Relations {
		items, err := relation.QueryTx(tx, data)
		if err != nil {
			logs.Error(err)
		}

		data[name] = items.Result
	}
}

/**
* getDetailsTx
* @param tx *Tx, data et.Json
* @return
**/
func (s *Ql) getDetailsTx(tx *Tx, data et.Json) {
	for name, relation := range s.Details {
		items, err := relation.QueryTx(tx, data)
		if err != nil {
			logs.Error(err)
		}

		data[name] = items.Result
	}
}

/**
* getCallsTx
* @param tx *Tx, data et.Json
* @return
**/
func (s *Ql) getCallsTx(tx *Tx, data et.Json) {
	for _, call := range s.Calcs {
		call(tx, data)
	}
}

/**
* queryTx
* @param tx *Tx
* @return (et.Items, error)
*
 */
func (s *Ql) queryTx(tx *Tx) (et.Items, error) {
	if s.db == nil {
		return et.Items{}, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	s.setTx(tx)
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
			s.getDetailsTx(tx, data)
		}(data)
	}
	wg.Wait()

	return result, nil
}

/**
* QueryTx
* @param tx *Tx, query et.Json
* @return et.Items, error
**/
func (s *Ql) QueryTx(tx *Tx, query et.Json) (et.Items, error) {
	s.setQuery(query)
	return s.queryTx(tx)
}

/**
* Query
* @param query et.Json
* @return et.Items, error
**/
func (s *Ql) Query(query et.Json) (et.Items, error) {
	return s.QueryTx(nil, query)
}

/**
* AllTx
* @param tx *Tx
* @return et.Items, error
**/
func (s *Ql) AllTx(tx *Tx) (et.Items, error) {
	return s.queryTx(tx)
}

/**
* All
* @return et.Items, error
**/
func (s *Ql) All() (et.Items, error) {
	return s.AllTx(nil)
}

/**
* LimitTx
* @param tx *Tx, page, rows int
* @return et.Items, error
**/
func (s *Ql) LimitTx(tx *Tx, page, rows int) (items et.Items, err error) {
	if rows > s.MaxRows {
		rows = s.MaxRows
	}
	s.Limits["page"] = page
	s.Limits["rows"] = rows
	return s.queryTx(tx)
}

/**
* Limit
* @param page, rows int
* @return et.Items, error
**/
func (s *Ql) Limit(page, rows int) (items et.Items, err error) {
	return s.LimitTx(nil, page, rows)
}

/**
* FirstTx
* @param tx *Tx, n int
* @return et.Items, error
**/
func (s *Ql) FirstTx(tx *Tx, n int) (et.Items, error) {
	return s.LimitTx(tx, 0, n)
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
* LastTx
* @param tx *Tx, n int
* @return et.Items, error
**/
func (s *Ql) LastTx(tx *Tx, n int) (et.Items, error) {
	return s.FirstTx(tx, n*-1)
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
* OneTx
* @param tx *Tx
* @return et.Item, error
**/
func (s *Ql) OneTx(tx *Tx) (et.Item, error) {
	result, err := s.LimitTx(tx, 0, 1)
	if err != nil {
		return et.Item{}, err
	}

	return result.First(), nil
}

/**
* One
* @return et.Item, error
**/
func (s *Ql) One() (et.Item, error) {
	return s.OneTx(nil)
}

/**
* RowsTx
* @param tx *Tx, limit int
* @return et.Items, error
**/
func (s *Ql) RowsTx(tx *Tx, val int) (et.Items, error) {
	page := s.Limits.Int("page")
	return s.LimitTx(tx, page, val)
}

/**
* Rows
* @param val int
* @return et.Items, error
**/
func (s *Ql) Rows(val int) (et.Items, error) {
	return s.RowsTx(nil, val)
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

	s.Exists = true
	result, err := s.queryTx(tx)
	if err != nil {
		return false, err
	}

	if result.Count == 0 {
		return false, nil
	}

	return result.Bool(0, "exists"), nil
}

/**
* ItExists
* @return bool, error
**/
func (s *Ql) ItExists() (bool, error) {
	return s.ItExistsTx(nil)
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

	s.Count = true
	result, err := s.queryTx(tx)
	if err != nil {
		return 0, err
	}

	if result.Count == 0 {
		return 0, nil
	}

	return result.Int(0, "all"), nil
}

/**
* Counted
* @return int, error
**/
func (s *Ql) Counted() (int, error) {
	return s.CountedTx(nil)
}

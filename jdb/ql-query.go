package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

/**
* FirstTx
* @param tx *Tx, n int
* @return et.Items, error
**/
func (s *Ql) FirstTx(tx *Tx, n int) (et.Items, error) {
	if s.Db == nil {
		return et.Items{}, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	s.setTx(tx)
	s.Limit = n
	s.prepare()
	result, err := s.Db.Select(s)
	if err != nil {
		return et.Items{}, err
	}

	for _, data := range result.Result {
		s.GetDetailsTx(tx, data)
	}

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
	if s.Db == nil {
		return false, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	s.setTx(tx)
	s.prepare()
	result, err := s.Db.Exists(s)
	if err != nil {
		return false, err
	}

	return result, nil
}

/**
* CountedTx
* @param tx *Tx
* @return int, error
**/
func (s *Ql) CountedTx(tx *Tx) (int, error) {
	if s.Db == nil {
		return 0, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	s.setTx(tx)
	s.prepare()
	result, err := s.Db.Count(s)
	if err != nil {
		return 0, err
	}

	return result, nil
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

/**
* QueryTx
* @param tx *Tx, params et.Json
* @return et.Json, error
**/
func (s *Ql) QueryTx(tx *Tx, params et.Json) (et.Json, error) {
	return s.queryTx(tx, params)
}

/**
* Query
* @param params et.Json
* @return et.Json, error
**/
func (s *Ql) Query(params et.Json) (et.Json, error) {
	return s.QueryTx(nil, params)
}

/**
* queryTx
* @param tx *Tx, params et.Json
* @return et.Items, error
**/
func (s *Ql) queryTx(tx *Tx, params et.Json) (et.Json, error) {
	if len(params) == 0 {
		return s.Help, nil
	}

	selects := params.Array("select")
	console.Pong()
	joins := params.ArrayJson("join")
	where := params.Json("where")
	groups := params.ArrayStr("group_by")
	havings := params.Json("having")
	orderBy := params.Json("order_by")
	page := params.Int("page")
	limit := params.ValInt(30, "limit")
	debug := params.Bool("debug")

	result, err := s.
		SetJoins(joins).
		SetWheres(where).
		SetGroupBy(groups...).
		SetHavings(havings).
		SetOrderBy(orderBy).
		SetSelects(selects...).
		SetDebug(debug).
		SetPage(page).
		SetLimitTx(tx, limit)
	if err != nil {
		return et.Json{}, err
	}

	return result, nil
}

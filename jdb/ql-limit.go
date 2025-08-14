package jdb

import (
	"github.com/cgalvisleon/et/config"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

/**
* ListTx
* @param tx *Tx, page, rows int
* @return et.List, error
**/
func (s *Ql) ListTx(tx *Tx, page, rows int) (et.List, error) {
	if s.Db == nil {
		return et.List{}, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	all, err := s.Db.Count(s)
	if err != nil {
		return et.List{}, err
	}

	s.Page(page)
	result, err := s.RowsTx(tx, rows)
	if err != nil {
		return et.List{}, err
	}

	return result.ToList(all, s.Sheet, s.Limit), nil
}

/**
* List
* @param page, rows int
* @return et.List, error
**/
func (s *Ql) List(page, rows int) (et.List, error) {
	return s.ListTx(s.tx, page, rows)
}

/**
* Page
* @param page int
* @return *Ql
**/
func (s *Ql) Page(val int) *Ql {
	s.Sheet = val
	return s
}

/**
* calcOffset
* @return *Ql
**/
func (s *Ql) calcOffset() *Ql {
	max := config.Int("DB_RECORD_LIMIT", 1000)
	if s.Limit > max {
		s.Limit = max
	}

	s.Offset = (s.Sheet - 1) * s.Limit
	if s.Offset < 0 {
		s.Offset = 0
	}

	return s
}

/**
* SetPage
* @param page int
* @return *Ql
**/
func (s *Ql) SetPage(page int) *Ql {
	s.Page(page)

	return s
}

/**
* SetLimitTx
* @param tx *Tx, limit int
* @return et.Json, error
**/
func (s *Ql) SetLimitTx(tx *Tx, limit int) (et.Json, error) {
	s.Limit = limit
	if s.Limit <= 0 {
		result, err := s.AllTx(tx)
		if err != nil {
			return nil, err
		}

		res := result.ToJson()
		if s.IsDebug {
			res["sql"] = s.Sql
		}

		return res, nil
	} else if s.Limit == 1 {
		result, err := s.OneTx(tx)
		if err != nil {
			return nil, err
		}

		res := result.ToJson()
		if s.IsDebug {
			res["sql"] = s.Sql
		}

		return res, nil
	} else {
		result, err := s.FirstTx(tx, s.Limit)
		if err != nil {
			return nil, err
		}

		res := result.ToJson()
		if s.IsDebug {
			res["sql"] = s.Sql
		}

		return res, nil
	}
}

/**
* getLimit
* @return interface{}
**/
func (s *Ql) getLimit() interface{} {
	if s.Sheet > 0 {
		return et.Json{
			"limit": s.Limit,
			"page":  s.Sheet,
		}
	}

	return s.Limit
}

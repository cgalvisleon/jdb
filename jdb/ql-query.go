package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

/**
* First
* @param n int
* @return et.Items, error
**/
func (s *Ql) First(n int) (et.Items, error) {
	if s.Db == nil {
		return et.Items{}, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	s.Limit = n
	s.prepare()
	result, err := s.Db.Select(s)
	if err != nil {
		return et.Items{}, err
	}

	for _, data := range result.Result {
		s.GetDetails(&data)
	}

	return result, nil
}

/**
* All
* @return et.Items, error
**/
func (s *Ql) All() (et.Items, error) {
	return s.First(0)
}

/**
* Last
* @param n int
* @return et.Items, error
**/
func (s *Ql) Last(n int) (et.Items, error) {
	if s.Db == nil {
		return et.Items{}, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	return s.First(n * -1)
}

/**
* One
* @return et.Item, error
**/
func (s *Ql) One() (et.Item, error) {
	result, err := s.First(1)
	if err != nil {
		return et.Item{}, err
	}

	if !result.Ok {
		return et.Item{Result: et.Json{}}, nil
	}

	return et.Item{
		Ok:     true,
		Result: result.Result[0],
	}, nil
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
* Rows
* @param limit int
* @return et.Items, error
**/
func (s *Ql) Rows(val int) (et.Items, error) {
	if s.Db == nil {
		return et.Items{}, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	return s.First(val)
}

/**
* List
* @param page, rows int
* @return et.List, error
**/
func (s *Ql) List(page, rows int) (et.List, error) {
	if s.Db == nil {
		return et.List{}, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	all, err := s.Db.Count(s)
	if err != nil {
		return et.List{}, err
	}

	s.Page(page)
	result, err := s.First(rows)
	if err != nil {
		return et.List{}, err
	}

	return result.ToList(all, s.Sheet, s.Limit), nil
}

/**
* Exist
* @return bool, error
**/
func (s *Ql) Exist() (bool, error) {
	if s.Db == nil {
		return false, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	s.prepare()
	result, err := s.Db.Exists(s)
	if err != nil {
		return false, err
	}

	return result, nil
}

/**
* Counted
* @return int, error
**/
func (s *Ql) Counted() (int, error) {
	if s.Db == nil {
		return 0, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	s.prepare()
	result, err := s.Db.Count(s)
	if err != nil {
		return 0, err
	}

	return result, nil
}

/**
* Query
* @param params et.Json
* @return interface{}, error
**/
func (s *Ql) Query(params et.Json) (interface{}, error) {
	if len(params) == 0 {
		return s.Help, nil
	}

	selects := params.Array("select")
	joins := params.ArrayJson("join")
	where := params.Json("where")
	groups := params.ArrayStr("group_by")
	havings := params.Json("having")
	orderBy := params.Json("order_by")
	page := params.Int("page")
	limit := params.ValInt(30, "limit")
	debug := params.Bool("debug")

	return s.setJoins(joins).
		setWheres(where).
		setGroupBy(groups...).
		setHavings(havings).
		setOrderBy(orderBy).
		setSelects(selects...).
		setDebug(debug).
		setPage(page).
		setLimit(limit)
}

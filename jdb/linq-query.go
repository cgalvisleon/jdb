package jdb

import (
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

/**
* First
* @param n int
* @return et.Items, error
**/
func (s *Linq) First(n int) (et.Items, error) {
	if s.Db == nil {
		return et.Items{}, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	limit := envar.GetInt(1000, "QUERY_LIMIT")
	if n > limit {
		n = limit
	}
	s.Limit = n
	s.Offset = (s.Sheet - 1) * s.Limit
	result, err := s.Db.Query(s)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* All
* @return et.Items, error
**/
func (s *Linq) All() (et.Items, error) {
	return s.First(0)
}

/**
* Last
* @param n int
* @return et.Items, error
**/
func (s *Linq) Last(n int) (et.Items, error) {
	if s.Db == nil {
		return et.Items{}, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	s.Limit = n
	result, err := s.Db.Last(s)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* One
* @return et.Item, error
**/
func (s *Linq) One() (et.Item, error) {
	result, err := s.First(1)
	if err != nil {
		return et.Item{}, err
	}

	if !result.Ok {
		return et.Item{}, nil
	}

	return et.Item{
		Ok:     true,
		Result: result.Result[0],
	}, nil
}

/**
* Offset
* @param offset int
* @return *Linq
**/
func (s *Linq) Page(val int) *Linq {
	s.Sheet = val
	s.Offset = (s.Sheet - 1) * s.Limit
	return s
}

/**
* Limit
* @param limit int
* @return *Linq
**/
func (s *Linq) Rows(val int) (et.Items, error) {
	if s.Db == nil {
		return et.Items{}, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	return s.First(val)
}

/**
* List
* @param page int
* @param rows int
* @return et.List, error
**/
func (s *Linq) List(page, rows int) (et.List, error) {
	if s.Db == nil {
		return et.List{}, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	all, err := s.Db.Count(s)
	if err != nil {
		return et.List{}, err
	}

	s.Sheet = page
	result, err := s.First(rows)
	if err != nil {
		return et.List{}, err
	}

	return result.ToList(all, s.Sheet, s.Limit), nil
}

/**
* Query
* @param query []string
* @return Linq
**/
func (s *Linq) Query(query et.Json) (et.Items, error) {
	selects := query.ArrayStr([]string{}, "select")
	joins := query.ArrayJson([]et.Json{}, "join")

	s.Select(selects...)
	s.SetJoins(joins)
	s.Db.Query(s)
	return et.Items{
		Ok: true,
		Result: []et.Json{{
			"select":   s.ListSelects(),
			"from":     s.ListForms(),
			"join":     s.ListJoins(),
			"where":    s.ListWheres(),
			"group_by": s.ListGroups(),
			"having":   s.ListHavings(),
			"order_by": s.ListOrders(),
			"limit": et.Json{
				"sheet":  s.Sheet,
				"offset": s.Offset,
			},
		}},
	}, nil
}

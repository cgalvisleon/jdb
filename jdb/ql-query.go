package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

/**
* GetDetails
* @return et.Json
**/
func (s *Ql) GetDetails(data et.Json) et.Json {

	return data
}

func (s *Ql) prepare() *Ql {
	for _, frm := range s.Froms {
		s.Selects = append(s.Selects, frm.Selects...)
	}

	if len(s.Selects) == 0 {
		frm := s.Froms[0]
		for _, col := range frm.Columns {
			if col.Hidden || slices.Contains([]TypeColumn{TpAtribute, TpGenerated}, col.TypeColumn) {
				continue
			}
			field := col.GetField()
			field.As = frm.As
			sel := &QlSelect{
				From:  frm,
				Field: field,
			}
			if col.TypeColumn == TpDetail {
				s.Details = append(s.Details, sel)
			} else {
				s.Selects = append(s.Selects, sel)
			}
		}
	}

	return s
}

/**
* First
* @param n int
* @return et.Items, error
**/
func (s *Ql) First(n int) (et.Items, error) {
	if s.Db == nil {
		return et.Items{}, mistake.New(MSG_DATABASE_NOT_FOUND)
	}

	s.setLimit(n)
	s.prepare()
	result, err := s.Db.Query(s)
	if err != nil {
		return et.Items{}, err
	}

	for i, item := range result.Result {
		result.Result[i] = s.GetDetails(item)
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
* Offset
* @param offset int
* @return *Ql
**/
func (s *Ql) Page(val int) *Ql {
	s.Sheet = val
	s.calcOffset()
	return s
}

/**
* Limit
* @param limit int
* @return *Ql
**/
func (s *Ql) Rows(val int) (et.Items, error) {
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
* Query
* @param params []string
* @return Ql
**/
func (s *Ql) Query(params et.Json) (et.Items, error) {
	joins := params.ArrayJson("join")
	where := params.ArrayJson("where")
	groups := params.ArrayStr("group_by")
	havings := params.ArrayJson("having")
	orders := params.ArrayJson("order_by")
	limit := params.ValInt(1000, "limit")
	page := params.ValInt(0, "page")

	if params["data"] != nil {
		data := params.ArrayStr("data")
		s.Data(data...)
	} else {
		selects := params.ArrayStr("select")
		s.Select(selects...)
	}
	s.setJoins(joins)
	s.setWheres(where)
	s.setGroupBy(groups...)
	s.setHavings(havings)
	s.setOrders(orders)
	s.setLimit(limit)
	s.setPage(page)
	s.Db.Query(s)
	return et.Items{
		Ok: true,
		Result: []et.Json{{
			"select":   s.listSelects(),
			"from":     s.listForms(),
			"join":     s.listJoins(),
			"where":    s.listWheres(),
			"group_by": s.listGroups(),
			"having":   s.listHavings(),
			"order_by": s.listOrders(),
			"limit":    s.listLimit(),
		}},
	}, nil
}

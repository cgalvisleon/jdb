package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

/**
* GetDetails
* @return et.Json
**/
func (s *Ql) GetDetails(data *et.Json) *et.Json {
	for _, field := range s.Details {
		col := field.Column
		if col == nil {
			continue
		}

		switch col.TypeColumn {
		case TpGenerated:
			if col.GeneratedFunction != nil {
				col.GeneratedFunction(col, data)
			}
		case TpRelatedTo:
			if col.Detail == nil {
				continue
			}
			if col.Detail.Fk == nil {
				continue
			}
			pkn := col.Detail.Fk.Name
			key := (*data)[pkn]
			if key == nil {
				continue
			}

			fkn := col.Detail.Key
			with := col.Detail.With
			limit := int(col.Detail.Limit)
			console.Debug("GetDetails:", " key:", key, " pkn:", pkn, " fkn:", fkn, " limit:", limit)
			if limit <= 0 {
				result, err := with.
					Where(fkn).Eq(key).
					Debug().
					All()
				if err != nil {
					continue
				}

				data.Set(col.Name, result.Result)
			} else {
				result, err := with.
					Where(fkn).Eq(key).
					Page(1).
					Debug().
					Rows(limit)
				if err != nil {
					continue
				}

				data.Set(col.Name, result.Result)
			}
		}
	}

	return data
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
* @param search et.Json
* @return Ql
**/
func (s *Ql) Query(search et.Json) (interface{}, error) {
	joins := search.ArrayJson("join")
	where := search.Json("where")
	groups := search.ArrayStr("group_by")
	havings := search.Json("having")
	orders := search.Json("order_by")
	page := search.Int("page")
	limit := search.ValInt(30, "limit")

	s.setJoins(joins).
		setWheres(where).
		setGroupBy(groups...).
		setHavings(havings).
		setOrders(orders)
	if search["data"] != nil {
		data := search.ArrayStr("data")
		s.Data(data...)
	} else {
		selects := search.ArrayStr("select")
		s.Select(selects...)
	}
	s.setPage(page)

	return s.setLimit(limit)
}

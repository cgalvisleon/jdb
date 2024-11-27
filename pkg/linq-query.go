package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
)

func (s *Linq) prebuild() {
	result := et.Json{}

	for _, val := range s.Froms {
		result.Set(val.Table, val.As)
	}

	for _, val := range s.Selects {
		result.Set(val.From.As, val.Field)
	}

	logs.Debug(result.ToString())
}

/**
* First
* @param n int
* @return et.Items, error
**/
func (s *Linq) First(n int) (et.Items, error) {
	s.prebuild()

	if s.Db.Driver == nil {
		return et.Items{}, logs.NewError(MSG_DRIVER_NOT_FOUND)
	}

	s.Limit = n
	result, err := (*s.Db.Driver).Query(s)
	if s.Show {
		logs.Debug(s.Describe().ToString())
	}
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
	if s.Db.Driver == nil {
		return et.Items{}, logs.NewError(MSG_DRIVER_NOT_FOUND)
	}

	s.Limit = n
	result, err := (*s.Db.Driver).Last(s)
	if s.Show {
		logs.Debug(s.Describe().ToString())
	}
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
	s.page = val
	s.Offset = s.Limit * (s.page - 1)
	return s
}

/**
* Limit
* @param limit int
* @return *Linq
**/
func (s *Linq) Rows(val int) (et.List, error) {
	if s.Db.Driver == nil {
		return et.List{}, logs.NewError(MSG_DRIVER_NOT_FOUND)
	}

	all, err := (*s.Db.Driver).Count(s)
	if s.Show {
		logs.Debug(s.Describe().ToString())
	}
	if err != nil {
		return et.List{}, err
	}

	s.Limit = val
	s.Offset = s.Limit * (s.page - 1)
	result, err := (*s.Db.Driver).Query(s)
	if s.Show {
		logs.Debug(s.Describe().ToString())
	}
	if err != nil {
		return et.List{}, err
	}

	return result.ToList(all, s.page, s.Limit), nil
}

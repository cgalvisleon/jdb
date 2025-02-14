package jdb

import "github.com/cgalvisleon/et/et"

type FilterTo interface {
	And(val interface{}) *QlFilter
	Or(val interface{}) *QlFilter
	Select(fields ...string) *Ql
	Data(fields ...string) *Ql
	Exec() (et.Items, error)
	One() (et.Item, error)
}

type QlFilter struct {
	main   FilterTo
	where  *QlWhere
	Wheres []*QlWhere
	Show   bool
}

func (s *QlFilter) setCondition(where et.Json) *QlFilter {
	if where["eq"] != nil {
		s.Eq(where["eq"])
	} else if where["neg"] != nil {
		s.Neg(where["neg"])
	} else if where["in"] != nil {
		s.In(where["in"])
	} else if where["like"] != nil {
		s.Like(where["like"])
	} else if where["more"] != nil {
		s.More(where["more"])
	} else if where["less"] != nil {
		s.Less(where["less"])
	} else if where["moreEq"] != nil {
		s.MoreEq(where["moreEq"])
	} else if where["lessEq"] != nil {
		s.LessEq(where["lessEq"])
	} else if where["search"] != nil {
		language := where.Str("language")
		s.Full(language, where["search"])
	} else if where["between"] != nil {
		s.Between(where["between"])
	} else if where["isNull"] != nil {
		s.IsNull()
	} else if where["notNull"] != nil {
		s.NotNull()
	}

	return s
}

/**
* AddValue
* @param val interface{}
* @return *QlFilter
**/
func (s *QlFilter) AddValue(val interface{}) FilterTo {
	appendValue := func(ql *Ql, value interface{}) {
		switch v := value.(type) {
		case string:
			field := ql.GetField(v, false)
			if field != nil {
				s.where.Values = append(s.where.Values, field)
			} else {
				s.where.Values = append(s.where.Values, value)
			}
		default:
			s.where.Values = append(s.where.Values, value)
		}
	}

	switch m := s.main.(type) {
	case *Ql:
		appendValue(m, val)
	case *Command:
		s.where.Values = append(s.where.Values, val)
	case *QlJoin:
		appendValue(m.Ql, val)
	case *QlHaving:
		appendValue(m.Ql, val)
	}

	s.Wheres = append(s.Wheres, s.where)
	return s.main
}

/**
* Eq
* @param val interface{}
* @return FilterTo
**/
func (s *QlFilter) Eq(val interface{}) FilterTo {
	s.where.Operator = Equal
	return s.AddValue(val)
}

/**
* Neg
* @param val interface{}
* @return FilterTo
**/
func (s *QlFilter) Neg(val interface{}) FilterTo {
	s.where.Operator = Neg
	return s.AddValue(val)
}

/**
* In
* @param val ...any
* @return FilterTo
**/
func (s *QlFilter) In(val ...any) FilterTo {
	s.where.Operator = In
	return s.AddValue(val)
}

/**
* Like
* @param val interface{}
* @return FilterTo
**/
func (s *QlFilter) Like(val interface{}) FilterTo {
	s.where.Operator = Like
	return s.AddValue(val)
}

/**
* More
* @param val interface{}
* @return FilterTo
**/
func (s *QlFilter) More(val interface{}) FilterTo {
	s.where.Operator = More
	return s.AddValue(val)
}

/**
* Less
* @param val interface{}
* @return FilterTo
**/
func (s *QlFilter) Less(val interface{}) FilterTo {
	s.where.Operator = Less
	return s.AddValue(val)
}

/**
* MoreEq
* @param val interface{}
* @return FilterTo
**/
func (s *QlFilter) MoreEq(val interface{}) FilterTo {
	s.where.Operator = MoreEq
	return s.AddValue(val)
}

/**
* LessEq
* @param val interface{}
* @return FilterTo
**/
func (s *QlFilter) LessEq(val interface{}) FilterTo {
	s.where.Operator = LessEq
	return s.AddValue(val)
}

/**
* Search
* @param val interface{}
* @return FilterTo
**/
func (s *QlFilter) Full(language string, val interface{}) FilterTo {
	s.where.Operator = Search
	s.where.Language = language
	return s.AddValue(val)
}

/**
* Between
* @param val1, val2 interface{}
* @return FilterTo
**/
func (s *QlFilter) Between(val interface{}) FilterTo {
	s.where.Operator = Between
	vals, ok := val.([]interface{})
	if !ok {
		return s.main
	}

	switch len(vals) {
	case 1:
		return s.AddValue(vals[0])
	case 2:
		s.AddValue(vals[0])
		return s.AddValue(vals[1])
	default:
		return s.main
	}
}

/**
* IsNull
* @return *QlFilter
**/
func (s *QlFilter) IsNull() *QlFilter {
	s.where.Operator = IsNull
	return s
}

/**
* NotNull
* @return *QlFilter
**/
func (s *QlFilter) NotNull() *QlFilter {
	s.where.Operator = NotNull
	return s
}

/**
* First
* @param n int
* @return et.Items, error
**/
func (s *QlFilter) First(n int) (et.Items, error) {
	return s.main.(*Ql).First(n)
}

/**
* All
* @return et.Items, error
**/
func (s *QlFilter) All() (et.Items, error) {
	return s.main.(*Ql).All()
}

/**
* Last
* @param n int
* @return et.Items, error
**/
func (s *QlFilter) Last(n int) (et.Items, error) {
	return s.main.(*Ql).Last(n)
}

/**
* One
* @return et.Item, error
**/
func (s *QlFilter) One() (et.Item, error) {
	return s.main.(*Ql).One()
}

/**
* Page
* @param val int
* @return *Ql
**/
func (s *QlFilter) Page(val int) *Ql {
	return s.main.(*Ql).Page(val)
}

/**
* Rows
* @param val int
* @return et.Items, error
**/
func (s *QlFilter) Rows(val int) (et.Items, error) {
	return s.main.(*Ql).Rows(val)
}

/**
* List
* @param page, rows int
* @return et.List, error
**/
func (s *QlFilter) List(page, rows int) (et.List, error) {
	return s.main.(*Ql).List(page, rows)
}

/**
* Query
* @param params et.Json
* @return (et.Items, error)
**/
func (s *QlFilter) Query(params et.Json) (et.Items, error) {
	return s.main.(*Ql).Query(params)
}

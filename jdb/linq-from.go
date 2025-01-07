package jdb

import "github.com/cgalvisleon/et/strs"

type LinqFrom struct {
	*Model
	As      string
	Selects []*LinqSelect
}

func From(m *Model) *Linq {
	result := &Linq{
		Db:     m.Db,
		Froms:  make([]*LinqFrom, 0),
		Joins:  make([]*LinqJoin, 0),
		Groups: make([]*LinqSelect, 0),
		Orders: make([]*LinqOrder, 0),
		Offset: 0,
		Limit:  0,
		Show:   false,
		Sheet:  0,
		index:  65,
	}
	result.LinqFilter = &LinqFilter{
		main:   result,
		Wheres: make([]*LinqWhere, 0),
	}
	result.Havings = &LinqHaving{
		Linq: result,
	}
	result.Havings.LinqFilter = &LinqFilter{
		main:   result.Havings,
		Wheres: make([]*LinqWhere, 0),
	}

	result.addFrom(m)

	return result
}

/**
* GetField
* @param name string
* @return *Field
**/
func (s *LinqFrom) GetField(name string) *Field {
	result := s.Model.GetField(name)
	if result != nil {
		result.As = s.As
	}

	return result
}

/**
* listForms
* @return []string
**/
func (s *Linq) listForms() []string {
	var result []string
	for _, from := range s.Froms {
		result = append(result, strs.Format(`%s, %s`, from.Table, from.As))
	}

	return result
}

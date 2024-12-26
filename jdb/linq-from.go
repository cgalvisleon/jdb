package jdb

import "github.com/cgalvisleon/et/strs"

type LinqFrom struct {
	*Model
	As      string
	Selects []*LinqSelect
}

func (s *LinqFrom) GetField(name string) *Field {
	result := NewField(name)
	if result == nil {
		return nil
	}

	col := s.GetColumn(result.Name)
	if col == nil && s.Integrity {
		return nil
	}

	if col == nil {
		col = s.DefineAtribute(name, TypeDataText)
	}

	result.Column = col
	result.Schema = s.Schema.Name
	result.Table = s.Name
	result.As = s.As
	result.Name = col.Field
	result.Atrib = col.Name

	return result
}

func From(m *Model) *Linq {
	result := &Linq{
		Db:      m.Db,
		Froms:   make([]*LinqFrom, 0),
		Joins:   make([]*LinqJoin, 0),
		Wheres:  make([]*LinqWhere, 0),
		Groups:  make([]*LinqSelect, 0),
		Havings: make([]*LinqWhere, 0),
		Orders:  make([]*LinqOrder, 0),
		Offset:  0,
		Limit:   0,
		Show:    false,
		Sheet:   1,
		index:   65,
	}

	result.addFrom(m)

	return result
}

/**
* ListForms
* @return []string
**/
func (s *Linq) ListForms() []string {
	var result []string
	if len(s.Froms) == 0 {
		return result
	}

	from := s.Froms[0]
	result = append(result, strs.Format(`%s, %s`, from.Table, from.As))

	return result
}

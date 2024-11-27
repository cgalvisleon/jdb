package jdb

func From(m *Model) *Linq {
	result := &Linq{
		Db:       m.Db,
		TypeLinq: TypeLinqSelect,
		Froms:    make([]*LinqFrom, 0),
		Joins:    make([]*LinqJoin, 0),
		Wheres:   make([]*LinqWhere, 0),
		GroupBys: make([]*LinqSelect, 0),
		Havings:  make([]*LinqWhere, 0),
		Selects:  make([]*LinqSelect, 0),
		Orders:   make([]*LinqOrder, 0),
		Offset:   0,
		Limit:    0,
		Show:     false,
		index:    65,
		page:     0,
	}

	result.addFrom(*m)

	return result
}

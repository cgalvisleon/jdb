package jdb

/**
* Order
* @param asc bool, fields ...string
* @return *Ql
**/
func (s *Ql) Order(asc bool, fields ...string) *Ql {
	if asc {
		s.OrdersBy["asc"] = fields
	} else {
		s.OrdersBy["desc"] = fields
	}

	return s
}

/**
* OrderBy
* @param fields ...string
* @return *Ql
**/
func (s *Ql) OrderBy(fields ...string) *Ql {
	s.OrdersBy["asc"] = fields
	return s
}

/**
* OrderDesc
* @param fields ...string
* @return *Ql
**/
func (s *Ql) OrderDesc(fields ...string) *Ql {
	s.OrdersBy["desc"] = fields
	return s
}

package jdb

/**
* BeforeInsert
* @param fn DataFunction
**/
func (s *Model) BeforeInsert(fn DataFunctionTx) *Model {
	s.beforeInsert = append(s.beforeInsert, fn)

	return s
}

/**
* BeforeUpdate
* @param fn DataFunction
**/
func (s *Model) BeforeUpdate(fn DataFunctionTx) *Model {
	s.beforeUpdate = append(s.beforeUpdate, fn)

	return s
}

/**
* BeforeDelete
* @param fn DataFunction
**/
func (s *Model) BeforeDelete(fn DataFunctionTx) *Model {
	s.beforeDelete = append(s.beforeDelete, fn)

	return s
}

/**
* BeforeInsertOrUpdate
* @param fn DataFunction
**/
func (s *Model) BeforeInsertOrUpdate(fn DataFunctionTx) *Model {
	s.beforeInsert = append(s.beforeInsert, fn)
	s.beforeUpdate = append(s.beforeUpdate, fn)

	return s
}

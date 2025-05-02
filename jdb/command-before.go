package jdb

/**
* BeforeInsert
* @param fn DataFunction
**/
func (s *Command) BeforeInsert(fn DataFunction) {
	s.beforeInsert = append(s.beforeInsert, fn)
}

/**
* BeforeUpdate
* @param fn DataFunction
**/
func (s *Command) BeforeUpdate(fn DataFunction) {
	s.beforeUpdate = append(s.beforeUpdate, fn)
}

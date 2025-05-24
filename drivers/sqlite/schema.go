package sqlite

/* Schema */
func (s *SqlLite) LoadSchema(name string) error
func (s *SqlLite) DropSchema(name string) error

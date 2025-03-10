package sqlite

import "github.com/cgalvisleon/jdb/jdb"

func (s *SqlLite) LoadTable(model *jdb.Model) (bool, error)
func (s *SqlLite) CreateModel(model *jdb.Model) error
func (s *SqlLite) DropModel(model *jdb.Model) error
func (s *SqlLite) SaveModel(model *jdb.Model) error

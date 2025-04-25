package sqlite

import "github.com/cgalvisleon/jdb/jdb"

func (s *SqlLite) LoadModel(model *jdb.Model) error
func (s *SqlLite) DropModel(model *jdb.Model) error
func (s *SqlLite) MutateModel(model *jdb.Model) error

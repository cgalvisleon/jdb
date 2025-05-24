package sqlite

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
)

func (s *SqlLite) Select(ql *jdb.Ql) (et.Items, error)
func (s *SqlLite) Count(ql *jdb.Ql) (int, error)
func (s *SqlLite) Exists(ql *jdb.Ql) (bool, error)

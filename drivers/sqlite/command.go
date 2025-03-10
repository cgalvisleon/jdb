package sqlite

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
)

func (s *SqlLite) Command(command *jdb.Command) (et.Items, error)

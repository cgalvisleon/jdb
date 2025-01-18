package base

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Base) Command(command *jdb.Command) (et.Items, error) {
	return et.Items{}, nil
}

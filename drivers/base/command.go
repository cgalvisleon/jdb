package base

import (
	"github.com/cgalvisl/jdb/jdb"
	"github.com/cgalvisleon/et/et"
)

func (s *Base) Current(command *jdb.Command) (et.Items, error) {
	return et.Items{}, nil
}

func (s *Base) Command(command *jdb.Command) (et.Item, error) {
	return et.Item{}, nil
}

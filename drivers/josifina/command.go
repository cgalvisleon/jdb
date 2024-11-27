package josefina

import (
	jdb "github.com/cgalvisl/jdb/pkg"
	"github.com/cgalvisleon/et/et"
)

func (s *Josefina) Current(command *jdb.Command) (et.Items, error) {
	return et.Items{}, nil
}

func (s *Josefina) Command(command *jdb.Command) (et.Item, error) {
	return et.Item{}, nil
}

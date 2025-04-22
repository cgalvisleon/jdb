package sqlite

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/jdb/jdb"
)

func (s *SqlLite) Command(command *jdb.Command) (et.Items, error) {
	return et.Items{}, mistake.New(MSG_COMMAND_NOT_FOUND)
}

func (s *SqlLite) Sync(command string, data et.Json) error {
	return mistake.New(MSG_FUNCION_NOT_FOUND)
}

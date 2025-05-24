package postgres

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

/**
* Sync
* @param command string
* @param data et.Json
* @return error
**/
func (s *Postgres) Sync(command string, data et.Json) error {
	return mistake.New(MSG_FUNCION_NOT_FOUND)
}

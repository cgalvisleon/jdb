package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* Command
* @param command *jdb.Command
* @return et.Items, error
**/
func (s *Postgres) Command(command *jdb.Command) (et.Items, error) {
	command.Sql = ""
	switch command.Command {
	case jdb.Insert:
		command.Sql = strs.Append(command.Sql, s.sqlInsert(command), "\n")
	case jdb.Update:
		command.Sql = strs.Append(command.Sql, s.sqlUpdate(command), "\n")
	case jdb.Delete:
		command.Sql = strs.Append(command.Sql, s.sqlDelete(command), "\n")
	}

	if command.IsDebug {
		console.Debug(command.Sql)
	}

	result, err := s.queryTx(command.Tx(), command.Sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* Sync
* @param command string
* @param data et.Json
* @return error
**/
func (s *Postgres) Sync(command string, data et.Json) error {
	return mistake.New(MSG_FUNCION_NOT_FOUND)
}

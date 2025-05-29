package mysql

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* Command
* @param command *jdb.Command
* @return et.Items, error
**/
func (s *Mysql) Command(command *jdb.Command) (et.Items, error) {
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

	result, err := jdb.QueryTx(command.Tx(), s.db, command.Sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

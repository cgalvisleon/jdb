package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) Command(command *jdb.Command) (et.Item, error) {
	command.Sql = ""
	switch command.Command {
	case jdb.Insert:
		command.Sql = strs.Append(command.Sql, s.sqlInsert(command), "\n")
	case jdb.Update:
		command.Sql = strs.Append(command.Sql, s.sqlUpdate(command), "\n")
	case jdb.Delete:
		command.Sql = strs.Append(command.Sql, s.sqlDelete(command), "\n")
	case jdb.Bulk:
		command.Sql = strs.Append(command.Sql, s.sqlBulk(command), "\n")
	}
	command.Sql = strs.Append(command.Sql, s.sqlReturn(command), "\n")

	if command.Show {
		console.Debug(command.Sql)
	}

	return et.Item{}, nil
}

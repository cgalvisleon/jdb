package oci

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) Command(command *jdb.Command) (et.Items, error) {
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

	if command.IsDebug {
		console.Debug(command.Sql)
	}

	if command.From.SourceField != nil {
		sourceField := command.From.SourceField.Name
		result, err := s.Data(sourceField, command.Sql)
		if err != nil {
			return et.Items{}, err
		}

		return result, nil
	}

	result, err := s.Query(command.Sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

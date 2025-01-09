package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlInsert(command *jdb.Command) string {
	result := "INSERT INTO %s(%s)\nVALUES (%s)"
	columns := ""
	values := ""
	console.Debug(command.Fields.ToString())
	console.Debug(command.Atribs.ToString())
	for key, val := range *command.New {
		column := strs.Uppcase(key)
		value := utility.Quote(val)

		columns = strs.Append(columns, column, ", ")
		values = strs.Append(values, strs.Format(`%v`, value), ", ")
	}

	result = strs.Format(result, command.Table, columns, values)

	return result
}

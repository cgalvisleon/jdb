package postgres

import (
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlInsert(command *jdb.Command) string {
	result := "INSERT INTO %s(%s)\nVALUES (%s)"
	columns := ""
	values := ""
	for key, val := range command.Fields {
		column := strs.Uppcase(key)
		value := utility.Quote(val)

		columns = strs.Append(columns, column, ", ")
		values = strs.Append(values, strs.Format(`%v`, value), ", ")
	}
	if len(command.Atribs) > 0 {
		columns = strs.Append(columns, command.From.SourceField.Up(), ", ")
		value := command.Atribs.ToString()

		values = strs.Append(values, strs.Format(`'%v'::jsonb`, value), ", ")
	}

	result = strs.Format(result, command.From.Table, columns, values)

	return result
}

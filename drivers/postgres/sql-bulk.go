package postgres

import (
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlBulk(command *jdb.Command) string {
	result := "INSERT INTO %s(%s)\nVALUES %s"
	columns := ""
	values := ""
	for i, data := range command.Origin {
		def := ""
		for key, val := range data {
			if i == 0 {
				column := strs.Uppcase(key)
				columns = strs.Append(columns, column, ", ")
			}

			value := utility.Quote(val)
			def = strs.Append(def, strs.Format(`%v`, value), ", ")
		}
		values = strs.Append(values, strs.Format(`(%v)`, def), ",\n")
	}

	result = strs.Format(result, command.From.Table, columns, values)

	return result
}

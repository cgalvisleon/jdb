package postgres

import (
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlInsert(command *jdb.Command) string {
	result := "INSERT INTO %s(%s)\nVALUE %s\nRETURNING jsonb_build_object() AS before,\n%s AS after;"
	columns := ""
	values := ""
	for i, value := range command.Values {
		record := ""
		for key, val := range value.Columns {
			if i == 0 {
				column := key
				columns = strs.Append(columns, column, ", ")
			}

			value := utility.Quote(val)
			def := strs.Format(`%v`, value)
			record = strs.Append(record, def, ", ")
		}
		if command.From.SourceField != nil && len(value.Atribs) > 0 {
			if i == 0 {
				column := command.From.SourceField.Name
				columns = strs.Append(columns, column, ", ")
			}

			value := value.Atribs.ToString()
			def := strs.Format(`'%v'::jsonb`, value)
			record = strs.Append(record, def, ", ")
		}
		def := strs.Format(`(%s)`, record)
		values = strs.Append(values, def, "\n")
	}

	objects := s.sqlJsonObject(command.From)
	return strs.Format(result, command.From.Table, columns, values, objects)
}

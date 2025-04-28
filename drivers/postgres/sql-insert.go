package postgres

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlInsert(command *jdb.Command) string {
	from := command.From
	columns := ""
	value := ""
	values := ""
	atribs := et.Json{}
	for _, val := range command.Values {
		for key, fld := range val {
			if fld.Column == from.SourceField || fld.Column == from.FullTextField {
				continue
			} else if fld.Column.TypeColumn == jdb.TpColumn {
				columns = strs.Append(columns, key, ", ")
				def := strs.Format(`%v`, fld.ValueQuoted())
				value = strs.Append(value, def, ", ")
			} else if fld.Column.TypeColumn == jdb.TpAtribute && from.SourceField != nil {
				atribs.Set(key, fld.Value)
			}
		}
		if from.SourceField != nil && len(atribs) > 0 {
			column := from.SourceField.Name
			columns = strs.Append(columns, column, ", ")

			def := strs.Format(`'%v'::jsonb`, atribs.ToString())
			value = strs.Append(value, def, ", ")
		}
		value = strs.Format(`(%s)`, value)
		values = strs.Append(values, value, ",\n")
	}

	objects := s.sqlJsonObject(from.GetFrom())
	returns := strs.Format("jsonb_build_object(\n'before', jsonb_build_object(),\n'after', (%s)) AS result;", objects)
	if len(command.Returns) > 0 {
		returns := ""
		for _, fld := range command.Returns {
			returns = strs.Append(returns, fld.Name, ", ")
		}
	}

	result := "INSERT INTO %s(%s)\nVALUES %s\nRETURNING\n%s"
	return strs.Format(result, table(from), columns, values, returns)
}

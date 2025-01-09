package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlUpdate(command *jdb.Command) string {
	result := "UPDATE %s SET\n%s\nWHERE %s"
	set := ""
	where := ""

	/*
		for key, val := range *command.New {
			column := strs.Uppcase(key)
			value := utility.Quote(val)

			if command.TypeSelect == jdb.Data && field == strs.Uppcase(SourceField.Upp()) {
				vals := strs.Uppcase(SourceField.Upp())
				atribs := c.new.Json(strs.Lowcase(field))

				for ak, av := range atribs {
					ak = strs.Lowcase(ak)
					av = et.Quote(av)

					vals = strs.Format(`jsonb_set(%s, '{%s}', '%v', true)`, vals, ak, av)
				}
				value = vals
			}

			fieldValue := strs.Format(`%s=%v`, column, value)
			set = strs.Append(set, fieldValue, ",\n")
		}*/

	result = strs.Format(result, command.Table, set, where)

	return result
}

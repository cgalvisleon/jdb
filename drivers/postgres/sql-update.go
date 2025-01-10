package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlUpdate(command *jdb.Command) string {
	result := "UPDATE %s SET\n%s\nWHERE %s"
	set := ""
	where := ""
	console.Debug(command.Fields.ToString())
	console.Debug(command.Atribs.ToString())

	for key, val := range command.Fields {
		column := strs.Uppcase(key)
		value := utility.Quote(val)

		if column == jdb.SourceField.Up() {
			continue
		}

		def := strs.Format(`%s=%v`, column, value)
		set = strs.Append(set, def, ",\n")
	}
	atribs := ""
	for key, val := range command.Atribs {
		atrib := strs.Uppcase(key)
		value := utility.Quote(val)

		if len(atribs) == 0 {
			atribs = jdb.SourceField.Up()
			atribs = strs.Format("jsonb_set(%s, '{%s}', %v, true)", atribs, atrib, value)
		} else {
			atribs = strs.Format("jsonb_set(\n%s, \n'{%s}', %v, true)", atribs, atrib, value)
		}
	}
	if len(atribs) > 0 {
		set = strs.Append(set, strs.Format(`%s=%s`, jdb.SourceField.Up(), atribs), ",\n")
	}

	frm := command.From
	where = whereFilters(command.Wheres)
	result = strs.Format(result, frm.Table, set, where)

	return result
}

package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlDelete(command *jdb.Command) string {
	result := "WITH updated_rows AS (\n\tSELECT\n\t\t%s AS _data\n\tFROM %s\n\tWHERE %s\n)\nDELETE FROM %s\nWHERE %s\nRETURNING\njsonb_build_object(\n'before', (%s),\n'after', jsonb_build_object(),\n'%s', %s) AS %s;"
	from := command.From
	where := whereFilters(command.Wheres)
	objects := s.sqlJsonObject(from)
	return strs.Format(result, objects, from.Table, where, from.Table, where, objects, jdb.SYSID, jdb.SYSID, command.Source)
}

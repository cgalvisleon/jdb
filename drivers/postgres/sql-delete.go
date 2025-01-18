package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlDelete(command *jdb.Command) string {
	result := "WITH updated_rows AS (\n\tSELECT\n\t\t%s AS _data\n\tFROM %s\n\tWHERE %s\n)\nDELETE FROM %s\nWHERE %s\nRETURNING (SELECT _data FROM updated_rows) AS before,\njsonb_build_object() AS after;"
	frm := command.From
	where := whereFilters(command.Wheres)
	return strs.Format(result, frm.Table, where)
}

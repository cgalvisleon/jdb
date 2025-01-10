package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlDelete(command *jdb.Command) string {
	frm := command.From
	where := whereFilters(command.Wheres)
	result := "DELETE FROM %s\nWHERE %s"
	result = strs.Format(result, frm.Table, where)

	return result
}

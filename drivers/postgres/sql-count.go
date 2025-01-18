package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) Count(ql *jdb.Ql) (int, error) {
	ql.Sql = ""
	ql.Sql = strs.Append(ql.Sql, "SELECT COUNT(*) AS Count", "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlFrom(ql.Froms), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlJoin(ql.Joins), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlWhere(ql.Wheres), "\n")

	if ql.Show {
		console.Debug(ql.Sql)
	}

	return 0, nil
}

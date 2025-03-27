package oci

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) Exists(ql *jdb.Ql) (bool, error) {
	ql.Sql = ""
	ql.Sql = strs.Append(ql.Sql, "SELECT 1", "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlFrom(ql.Froms), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlJoin(ql.Joins), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlWhere(ql.QlWhere), "\n")

	if len(ql.Sql) > 0 {
		ql.Sql = strs.Format("SELECT EXISTS (%s);", ql.Sql)
	}

	if ql.IsDebug {
		console.Debug(ql.Sql)
	}

	result, err := s.Query(ql.Sql)
	if err != nil {
		return false, err
	}

	if result.Count == 0 {
		return false, nil
	}

	return result.Bool(0, "exists"), nil
}

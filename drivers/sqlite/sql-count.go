package sqlite

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* Count
* @param ql *jdb.Ql
* @return int, error
**/
func (s *SqlLite) Count(ql *jdb.Ql) (int, error) {
	ql.Sql = ""
	ql.Sql = strs.Append(ql.Sql, "SELECT COUNT(*) AS Count", "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlFrom(ql.Froms), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlJoin(ql.Joins), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlWhere(ql.QlWhere), "\n")

	if ql.IsDebug {
		console.Debug(ql.Sql)
	}

	result, err := jdb.Query(s.db, ql.Sql)
	if err != nil {
		return 0, err
	}

	if result.Count == 0 {
		return 0, nil
	}

	return result.Int(0, "count"), nil
}

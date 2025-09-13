package sqlite

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* Exists
* @param ql *jdb.Ql
* @return bool, error
**/
func (s *SqlLite) Exists(ql *jdb.Ql) (bool, error) {
	ql.Sql = ""
	ql.Sql = strs.Append(ql.Sql, "SELECT 1", "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlFrom(ql.Froms), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlWhere(ql.QlWhere), "\n")

	if len(ql.Sql) > 0 {
		ql.Sql = strs.Format(`SELECT EXISTS (%s) AS "exists";`, ql.Sql)
	}

	if ql.IsDebug {
		console.Debug(ql.Sql)
	}

	item, err := jdb.Query(s.jdb, ql.Sql)
	if err != nil {
		return false, err
	}

	result := item.Int(0, "exists") > 0

	return result, nil
}

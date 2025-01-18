package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlHaving(ql *jdb.Ql) string {
	wheres := ql.Havings.Wheres
	where := whereFilters(wheres)
	result := "HAVING %s"
	result = strs.Format(result, where)

	return result
}

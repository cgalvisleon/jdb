package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlHaving(linq *jdb.Ql) string {
	wheres := linq.Havings.Wheres
	where := whereFilters(wheres)
	result := "HAVING %s"
	result = strs.Format(result, where)

	return result
}

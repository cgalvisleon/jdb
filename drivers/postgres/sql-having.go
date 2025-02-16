package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlHaving(ql *jdb.Ql) string {
	result := ""
	havings := ql.Havings
	where := whereConditions(havings.QlWhere)
	if where != "" {
		result = strs.Format("HAVING %s", where)
	}

	return result
}

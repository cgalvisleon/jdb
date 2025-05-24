package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* sqlHaving
* @param ql *jdb.Ql
* @return string
**/
func (s *Postgres) sqlHaving(ql *jdb.Ql) string {
	result := ""
	havings := ql.Havings
	where := whereConditions(havings.QlWhere)
	if where == "" {
		return result
	}

	result = strs.Format("HAVING %s", where)

	return result
}

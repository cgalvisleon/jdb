package postgres

import (
	"fmt"

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

	result = fmt.Sprintf("HAVING %s", where)

	return result
}

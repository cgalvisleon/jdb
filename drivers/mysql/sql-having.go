package mysql

import (
	"fmt"

	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* sqlHaving
* @param ql *jdb.Ql
* @return string
**/
func (s *Mysql) sqlHaving(ql *jdb.Ql) string {
	result := ""
	havings := ql.Havings
	where := whereConditions(havings.QlWhere)
	if where == "" {
		return result
	}

	result = fmt.Sprintf("HAVING %s", where)

	return result
}

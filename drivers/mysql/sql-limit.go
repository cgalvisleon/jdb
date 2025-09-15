package mysql

import (
	"fmt"

	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* sqlLimit
* @param ql *jdb.Ql
* @return string
**/
func (s *Mysql) sqlLimit(ql *jdb.Ql) string {
	result := ""
	if ql.Sheet > 0 {
		result = fmt.Sprintf(`LIMIT %d OFFSET %d`, ql.Limit, ql.Offset)
	} else if ql.Limit > 0 {
		result = fmt.Sprintf(`LIMIT %d`, ql.Limit)
	}

	return result
}

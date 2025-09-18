package sqlite

import (
	"fmt"

	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* sqlGroupBy
* @param ql *jdb.Ql
* @return string
**/
func (s *SqlLite) sqlGroupBy(ql *jdb.Ql) string {
	result := ""
	columns := s.sqlColumns(ql.Groups)
	if len(columns) == 0 {
		return result
	}

	result = fmt.Sprintf("GROUP BY %s", columns)

	return result
}

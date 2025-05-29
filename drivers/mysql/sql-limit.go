package mysql

import (
	"github.com/cgalvisleon/et/strs"
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
		result = strs.Format(`LIMIT %d OFFSET %d`, ql.Limit, ql.Offset)
	} else if ql.Limit > 0 {
		result = strs.Format(`LIMIT %d`, ql.Limit)
	}

	return result
}

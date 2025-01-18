package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlGroupBy(ql *jdb.Ql) string {
	result := "GROUP BY %s"
	columns := s.sqlColumns(ql.Groups)
	result = strs.Format(result, columns)

	return result
}

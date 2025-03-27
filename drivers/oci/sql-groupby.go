package oci

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlGroupBy(ql *jdb.Ql) string {
	result := ""
	columns := s.sqlColumns(ql.Groups)
	if len(columns) == 0 {
		return result
	}

	result = strs.Format("GROUP BY %s", columns)

	return result
}

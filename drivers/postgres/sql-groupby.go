package postgres

import jdb "github.com/cgalvisleon/jdb/jdb"

func (s *Postgres) sqlGroupBy(linq *jdb.Ql) string {
	result := "GROUP BY %s"
	result = s.sqlColumns(nil, linq.TypeSelect, linq.Groups, nil)

	return result
}

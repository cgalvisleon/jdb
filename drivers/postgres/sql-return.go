package postgres

import (
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlJsonObject(from *jdb.QlFrom) string {
	var selects = []*jdb.Field{}
	for _, col := range from.Columns {
		field := col.GetField()
		selects = append(selects, field)
	}

	return s.sqlObject(selects)
}

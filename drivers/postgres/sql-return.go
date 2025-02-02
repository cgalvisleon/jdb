package postgres

import (
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlJsonObject(from *jdb.QlFrom) string {
	var selects = []*jdb.QlSelect{}
	from.SetSelectByColumns(&selects, nil)

	return s.sqlObject(selects)
}

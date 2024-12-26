package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) queryFrom(froms []*jdb.LinqFrom) string {
	if len(froms) == 0 {
		return ""
	}

	from := froms[0]
	result := strs.Format("FROM %s AS %s", from.Table, from.As)

	return result
}

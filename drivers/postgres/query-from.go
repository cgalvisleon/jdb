package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) tableAs(from *jdb.LinqFrom) string {
	if from == nil {
		return ""
	}

	return strs.Append(from.Table, from.As, " AS ")
}

func (s *Postgres) queryFrom(froms []*jdb.LinqFrom) string {
	if len(froms) == 0 {
		return ""
	}

	from := froms[0]
	def := s.tableAs(from)
	result := strs.Format("FROM %s", def)

	return result
}

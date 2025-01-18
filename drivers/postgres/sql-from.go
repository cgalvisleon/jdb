package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) tableAs(from *jdb.QlFrom) string {
	if from == nil {
		return ""
	}

	return strs.Append(from.Table, from.As, " AS ")
}

func (s *Postgres) sqlFrom(froms []*jdb.QlFrom) string {
	if len(froms) == 0 {
		return ""
	}

	from := froms[0]
	def := s.tableAs(from)
	result := strs.Format("FROM %s", def)

	return result
}

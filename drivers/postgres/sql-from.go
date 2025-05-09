package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlFrom(froms *jdb.QlFroms) string {
	if len(froms.Froms) == 0 {
		return ""
	}

	from := froms.Froms[0]
	def := s.tableAs(from)
	result := strs.Format("FROM %s", def)

	return result
}

func (s *Postgres) tableAs(from *jdb.QlFrom) string {
	if from == nil {
		return ""
	}

	return strs.Append(table(from.Model), from.As, " AS ")
}

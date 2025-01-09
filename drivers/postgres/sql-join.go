package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlJoin(joins []*jdb.LinqJoin) string {
	result := ""
	for _, join := range joins {
		def := s.tableAs(join.From)
		def = strs.Append(def, whereFilters(join.Wheres), " ON ")
		switch join.TypeJoin {
		case jdb.JoinInner:
			def = strs.Append(`INNER JOIN`, def, " ")
			result = strs.Append(result, def, "\n")
		case jdb.JoinLeft:
			def = strs.Append(`LEFT JOIN`, def, " ")
			result = strs.Append(result, def, "\n")
		case jdb.JoinRight:
			def = strs.Append(`RIGHT JOIN`, def, " ")
			result = strs.Append(result, def, "\n")
		case jdb.JoinFull:
			def = strs.Append(`FULL JOIN`, def, " ")
			result = strs.Append(result, def, "\n")
		}
	}

	return result
}

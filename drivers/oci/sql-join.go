package oci

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Oracle) sqlJoin(joins []*jdb.QlJoin) string {
	result := ""
	for _, join := range joins {
		def := s.tableAs(join.With)
		def = strs.Append(def, whereConditions(join.QlWhere), " ON ")
		switch join.TypeJoin {
		case jdb.InnerJoin:
			def = strs.Append(`INNER JOIN`, def, " ")
			result = strs.Append(result, def, "\n")
		case jdb.LeftJoin:
			def = strs.Append(`LEFT JOIN`, def, " ")
			result = strs.Append(result, def, "\n")
		case jdb.RightJoin:
			def = strs.Append(`RIGHT JOIN`, def, " ")
			result = strs.Append(result, def, "\n")
		case jdb.FullJoin:
			def = strs.Append(`FULL JOIN`, def, " ")
			result = strs.Append(result, def, "\n")
		}
	}

	return result
}

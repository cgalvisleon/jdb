package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlOrderBy(ql *jdb.Ql) string {
	result := ""
	for _, ord := range ql.Orders {
		def := selectField(ord.Field)
		if ord.Sorted {
			def = strs.Append(def, "ASC", " ")
		} else {
			def = strs.Append(def, "DESC", " ")
		}
		result = strs.Append(result, def, ",\n")
	}
	if len(result) != 0 {
		result = strs.Append("ORDER BY", result, "\n")
	}

	return result
}

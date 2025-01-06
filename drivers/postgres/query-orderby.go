package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) queryOrderBy(linq *jdb.Linq) string {
	result := ""
	for _, ord := range linq.Orders {
		def := selectField(ord.Field)
		if ord.Sorted {
			def = strs.Append(def, "ASC", " ")
		} else {
			def = strs.Append(def, "DESC", " ")
		}
		result = strs.Append(result, def, ",\n")
	}
	result = strs.Append("ORDER BY", result, "\n")

	return result
}

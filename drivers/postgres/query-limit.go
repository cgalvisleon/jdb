package postgres

import (
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) queryLimit(linq *jdb.Linq) string {
	result := ""
	if linq.Sheet > 0 {
		result = strs.Format(`LIMIT %d OFFSET %d`, linq.Limit, linq.Offset)
	} else if linq.Limit > 0 {
		result = strs.Format(`LIMIT %d`, linq.Limit)
	} else {
		limit := envar.GetInt(1000, "QUERY_LIMIT")
		result = strs.Format(`LIMIT %d`, limit)
	}

	return result
}

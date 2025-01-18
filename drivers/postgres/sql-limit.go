package postgres

import (
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlLimit(ql *jdb.Ql) string {
	result := ""
	if ql.Sheet > 0 {
		result = strs.Format(`LIMIT %d OFFSET %d`, ql.Limit, ql.Offset)
	} else if ql.Limit > 0 {
		result = strs.Format(`LIMIT %d`, ql.Limit)
	} else {
		limit := envar.GetInt(1000, "QUERY_LIMIT")
		result = strs.Format(`LIMIT %d`, limit)
	}

	return result
}

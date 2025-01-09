package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlReturn(command *jdb.Command) string {
	var selects = []*jdb.LinqSelect{}
	var orders = []*jdb.LinqOrder{}

	selects = append(selects, command.Returns...)
	if len(selects) == 0 {
		return ""
	}

	result := s.sqlColumns(command.TypeSelect, selects, orders)
	result = strs.Append("RETURNING", result, "\n")

	return result
}

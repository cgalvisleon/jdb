package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlReturn(command *jdb.Command) string {
	var selects = []*jdb.LinqSelect{}
	var orders = []*jdb.LinqOrder{}

	frm := command.From
	for _, col := range frm.Columns {
		if col.TypeColumn != jdb.TpColumn {
			continue
		}
		field := col.GetField()
		field.As = frm.As
		selects = append(selects, &jdb.LinqSelect{
			From:  frm,
			Field: field,
		})
	}

	result := s.sqlColumns(frm, command.TypeSelect, selects, orders)
	result = strs.Append("RETURNING", result, "\n")

	return result
}

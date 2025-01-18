package postgres

import (
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlJsonObject(from *jdb.QlFrom) string {
	var selects = []*jdb.QlSelect{}
	for _, col := range from.Columns {
		if col.TypeColumn != jdb.TpColumn {
			continue
		}
		field := col.GetField()
		field.As = from.As
		selects = append(selects, &jdb.QlSelect{
			From:  from,
			Field: field,
		})
	}

	return s.sqlData(from, selects)
}

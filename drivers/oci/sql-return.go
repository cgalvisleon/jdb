package oci

import (
	"slices"

	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Oracle) sqlJsonObject(from *jdb.QlFrom) string {
	var selects = []*jdb.Field{}
	for _, col := range from.Columns {
		field := col.GetField()
		if field == nil {
			continue
		}
		if field.Column == nil {
			continue
		}
		if slices.Contains([]jdb.TypeColumn{jdb.TpColumn}, field.Column.TypeColumn) {
			field.As = from.As
			selects = append(selects, field)
		}
	}

	return s.sqlObject(selects)
}

package oci

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Oracle) sqlFrom(froms *jdb.QlFroms) string {
	if len(froms.Froms) == 0 {
		return ""
	}

	from := froms.Froms[0]
	def := s.tableAs(from)
	result := strs.Format("FROM %s", def)

	return result
}

func (s *Oracle) tableAs(from *jdb.QlFrom) string {
	if from == nil {
		return ""
	}

	return strs.Append(from.Table, from.As, " AS ")
}

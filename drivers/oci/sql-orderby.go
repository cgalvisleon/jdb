package oci

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Oracle) sqlOrderBy(ql *jdb.Ql) string {
	result := ""
	for _, fld := range ql.Orders.Asc {
		def := asField(*fld)
		def = strs.Append(def, "ASC", " ")
		result = strs.Append(result, def, ",\n")
	}
	for _, fld := range ql.Orders.Desc {
		def := asField(*fld)
		def = strs.Append(def, "DESC", " ")
		result = strs.Append(result, def, ",\n")
	}

	if len(result) != 0 {
		result = strs.Append("ORDER BY", result, "\n")
	}

	return result
}

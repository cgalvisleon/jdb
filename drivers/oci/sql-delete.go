package oci

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlDelete(command *jdb.Command) string {
	from := command.From
	where := whereConditions(command.QlWhere)
	objects := s.sqlJsonObject(from.GetFrom())

	result := "WITH deleted_rows AS (\nSELECT\nctid,\n%s AS old_data\nFROM %s\nWHERE %s\n)\nDELETE FROM %s AS oc\nUSING deleted_rows dr\nWHERE oc.ctid = dr.ctid\nRETURNING\njsonb_build_object(\n'before', (dr.old_data),\n'after', jsonb_build_object()) AS result;"
	return strs.Format(result, objects, from.Table, where, from.Table)
}

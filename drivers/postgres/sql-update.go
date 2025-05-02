package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlUpdate(command *jdb.Command) string {
	from := command.From
	set := ""
	atribs := ""
	where := ""
	for _, value := range command.Values {
		for key, field := range value {
			if field.Column.TypeColumn == jdb.TpColumn {
				val := field.ValueQuoted()
				def := strs.Format(`%s = %v`, key, val)
				set = strs.Append(set, def, ",\n")
			} else if field.Column.TypeColumn == jdb.TpAtribute && from.SourceField != nil {
				val := JsonQuote(field.Value)
				if len(atribs) == 0 {
					atribs = from.SourceField.Name
					atribs = strs.Format("jsonb_set(%s, '{%s}', %v::jsonb, true)", atribs, key, val)
				} else {
					atribs = strs.Format("jsonb_set(\n%s, \n'{%s}', %v::jsonb, true)", atribs, key, val)
				}
			}
		}
		if len(atribs) > 0 {
			def := strs.Format(`%s = %v`, from.SourceField.Name, atribs)
			set = strs.Append(set, def, ",\n")
		}
	}

	where = whereConditions(command.QlWhere)
	objects := s.sqlJsonObject(from.GetFrom())
	returns := strs.Format("jsonb_build_object(\n'before', (ur.old_data),\n'after', (%s)) AS result;", objects)
	if len(command.Returns) > 0 {
		returns = ""
		for _, fld := range command.Returns {
			returns = strs.Append(returns, fld.Name, ", ")
		}
	}

	result := "WITH updated_rows AS (\nSELECT\noc.ctid,\n%s AS old_data\nFROM %s AS oc\nWHERE %s)\nUPDATE %s AS oc SET\n%s\nFROM updated_rows ur\nWHERE oc.ctid = ur.ctid\nRETURNING\n%s;"
	return strs.Format(result, objects, table(from), where, table(from), set, returns)
}

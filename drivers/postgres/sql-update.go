package postgres

import (
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlUpdate(command *jdb.Command) string {
	result := "WITH updated_rows AS (\n\tSELECT\n\t\t%s AS _data\n\tFROM %s\n\tWHERE %s\n)\nUPDATE %s SET\n%s\nWHERE %sRETURNING\n\t(SELECT _data FROM updated_rows) AS before,\n%s AS after;"
	from := command.From
	set := ""
	atribs := ""
	where := ""
	for _, value := range command.Values {
		for key, val := range value.Columns {
			val := utility.Quote(val)
			if key == string(jdb.SourceField) {
				continue
			}

			def := strs.Format(`%s = %v`, key, val)
			set = strs.Append(set, def, ",\n")
		}
		for key, val := range value.Atribs {
			val := utility.Quote(val)
			if len(atribs) == 0 {
				atribs = string(jdb.SourceField)
				atribs = strs.Format("jsonb_set(%s, '{%s}', %v, true)", atribs, key, val)
			} else {
				atribs = strs.Format("jsonb_set(\n%s, \n'{%s}', %v, true)", atribs, key, val)
			}
		}
		if len(atribs) > 0 {
			def := strs.Format(`%s = %v`, string(jdb.SourceField), atribs)
			set = strs.Append(set, def, ",\n")
		}
		if from.FullTextField != nil {
			tsvector := ""
			for _, key := range from.FullTextField.FullText {
				if value.Data[key] != nil {
					val := utility.Quote(value.Data[key])
					tsvector = strs.Append(tsvector, strs.Format(`%v`, val), " || '' || ")
				}
			}
			def := strs.Format(`to_tsvector('%s', %s)`, from.FullTextField.Language, value)
			def = strs.Format(`%s = %v`, from.FullTextField.Name, def)
			set = strs.Append(set, def, ",\n")
		}
	}

	where = whereFilters(command.Wheres)
	objects := s.sqlJsonObject(from)
	result = strs.Format(result, objects, from.Table, where, from.Table, set, where, objects)

	return result
}

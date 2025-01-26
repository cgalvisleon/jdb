package postgres

import (
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlUpdate(command *jdb.Command) string {
	result := "WITH updated_rows AS (\nSELECT\n%s AS _data\nFROM %s\nWHERE %s)\nUPDATE %s SET\n%s\nWHERE %sRETURNING\njsonb_build_object(\n'before', (SELECT _data FROM updated_rows),\n'after', (%s),\n'%s', %s) AS %s;"
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
			val := JsonQuote(val)
			if len(atribs) == 0 {
				atribs = string(jdb.SourceField)
				atribs = strs.Format("jsonb_set(%s, '{%s}', %v::jsonb, true)", atribs, key, val)
			} else {
				atribs = strs.Format("jsonb_set(\n%s, \n'{%s}', %v::jsonb, true)", atribs, key, val)
			}
		}
		if len(atribs) > 0 {
			def := strs.Format(`%s = %v`, string(jdb.SourceField), atribs)
			set = strs.Append(set, def, ",\n")
		}
		if from.FullTextField != nil {
			tsvector := ""
			for _, key := range from.FullTextField.FullText {
				v := value.Data[key]
				if v != nil {
					val := utility.Quote(v)
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
	result = strs.Format(result, objects, from.Table, where, from.Table, set, where, objects, jdb.SYSID, jdb.SYSID, jdb.SOURCE)

	return result
}

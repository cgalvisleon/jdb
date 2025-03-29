package oci

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Oracle) sqlInsert(command *jdb.Command) string {
	from := command.From
	columns := ""
	values := ""
	atribs := et.Json{}
	for _, value := range command.Values {
		for key, fld := range value {
			if fld.Column == from.SourceField || fld.Column == from.FullTextField {
				continue
			} else if fld.Column.TypeColumn == jdb.TpColumn {
				columns = strs.Append(columns, key, ", ")
				def := strs.Format(`%v`, utility.Quote(fld.Value))
				values = strs.Append(values, def, ", ")
			} else if fld.Column.TypeColumn == jdb.TpAtribute && from.SourceField != nil {
				atribs.Set(key, fld.Value)
			}
		}
		if from.SourceField != nil && len(atribs) > 0 {
			column := from.SourceField.Name
			columns = strs.Append(columns, column, ", ")

			def := strs.Format(`'%v'::jsonb`, atribs.ToString())
			values = strs.Append(values, def, ", ")
		}
	}

	objects := s.sqlJsonObject(from.GetFrom())
	result := "INSERT INTO %s(%s)\nVALUES (%s)\nRETURNING\njsonb_build_object(\n'before', jsonb_build_object(),\n'after', (%s)) AS result;"
	return strs.Format(result, from.Table, columns, values, objects)
}

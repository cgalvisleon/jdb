package postgres

import (
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlInsert(command *jdb.Command) string {
	result := "INSERT INTO %s(%s)\nVALUES %s\nRETURNING\njsonb_build_object(\n'before', jsonb_build_object(),\n'after', (%s),\n'%s', %s) AS %s;"
	from := command.From
	columns := ""
	values := ""
	for i, value := range command.Values {
		record := ""
		for key, val := range value.Columns {
			if i == 0 {
				column := key
				columns = strs.Append(columns, column, ", ")
			}

			value := utility.Quote(val)
			def := strs.Format(`%v`, value)
			record = strs.Append(record, def, ", ")
		}
		if from.SourceField != nil && len(value.Atribs) > 0 {
			if i == 0 {
				column := from.SourceField.Name
				columns = strs.Append(columns, column, ", ")
			}

			value := value.Atribs.ToString()
			def := strs.Format(`'%v'::jsonb`, value)
			record = strs.Append(record, def, ", ")
		}
		if from.FullTextField != nil {
			if i == 0 {
				column := from.FullTextField.Name
				columns = strs.Append(columns, column, ", ")
			}

			tsvector := ""
			for _, key := range from.FullTextField.FullText {
				if value.Data[key] != nil {
					val := utility.Quote(value.Data[key])
					tsvector = strs.Append(tsvector, strs.Format(`%v`, val), " || '' || ")
				}
			}
			def := strs.Format(`to_tsvector('%s', '%s')`, from.FullTextField.Language, value)
			record = strs.Append(record, def, ", ")
		}
		def := strs.Format(`(%s)`, record)
		values = strs.Append(values, def, "\n")
	}

	objects := s.sqlJsonObject(from)
	return strs.Format(result, from.Table, columns, values, objects, jdb.SYSID, jdb.SYSID, command.Source)
}

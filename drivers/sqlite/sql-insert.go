package sqlite

import (
	"fmt"
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* sqlInsert
* @param command *jdb.Command
* @return string
**/
func (s *SqlLite) sqlInsert(command *jdb.Command) string {
	from := command.From
	columns := utility.NewList()
	value := ""
	values := ""
	atribs := et.Json{}
	for _, val := range command.Values {
		for key, field := range val {
			switch field.Column.TypeColumn {
			case jdb.TpColumn:
				columns.Add(key)
				def := fmt.Sprintf(`%v`, field.ValueQuoted())
				value = strs.Append(value, def, ", ")
			case jdb.TpAtribute:
				atribs.Set(key, field.Value)
			}
		}

		if from.SourceField != nil && len(atribs) > 0 {
			column := from.SourceField.Name
			columns.Add(column)

			def := fmt.Sprintf(`'%s'`, atribs.ToString())
			value = strs.Append(value, def, ", ")
		}

		value = fmt.Sprintf(`(%s)`, value)
		values = strs.Append(values, value, ",\n")
	}

	objects := s.sqlObject(from.GetFrom())
	returns := fmt.Sprintf("%s AS result", objects)
	if len(command.Returns) > 0 {
		returns := ""
		for _, fld := range command.Returns {
			returns = strs.Append(returns, fld.Name, ", ")
		}
	}

	columnNames := []string{}
	for _, column := range columns {
		columnNames = append(columnNames, column.(string))
	}

	result := "INSERT INTO %s(%s)\nVALUES %s\nRETURNING\n%s;"
	return fmt.Sprintf(result, tableName(from), strings.Join(columnNames, ", "), values, returns)
}

package sqlite

import (
	"fmt"

	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* sqlUpdate
* @param command *jdb.Command
* @return string
**/
func (s *SqlLite) sqlUpdate(command *jdb.Command) string {
	from := command.From
	set := ""
	atribs := ""
	where := ""
	for _, value := range command.Values {
		for key, field := range value {
			if field.Column.TypeColumn == jdb.TpColumn {
				val := field.ValueQuoted()
				def := fmt.Sprintf(`%s = %v`, key, val)
				set = strs.Append(set, def, ",\n")
			} else if field.Column.TypeColumn == jdb.TpAtribute && from.SourceField != nil {
				val := jdb.JsonQuote(field.Value)
				if len(atribs) == 0 {
					atribs = from.SourceField.Name
					atribs = fmt.Sprintf("json_set(%s, '$.%s', %v)", atribs, key, val)
				} else {
					atribs = fmt.Sprintf("json_set(\n%s, \n'$.%s', %v)", atribs, key, val)
				}
			}
		}
		if len(atribs) > 0 {
			def := fmt.Sprintf(`%s = %v`, from.SourceField.Name, atribs)
			set = strs.Append(set, def, ",\n")
		}
	}

	where = whereConditions(command.QlWhere)
	objects := s.sqlObject(from.GetFrom())
	returns := fmt.Sprintf("%s AS result", objects)
	if len(command.Returns) > 0 {
		returns = ""
		for _, fld := range command.Returns {
			returns = strs.Append(returns, fld.Name, ", ")
		}
	}

	result := "UPDATE %s SET\n%s\nWHERE %s\nRETURNING\n%s;"
	return fmt.Sprintf(result, tableName(from), set, where, returns)
}

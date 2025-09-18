package sqlite

import (
	"fmt"

	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* SqlDelete
* @param command *jdb.Command
* @return string
**/
func (s *SqlLite) sqlDelete(command *jdb.Command) string {
	from := command.GetFrom()
	if from == nil {
		return ""
	}

	where := whereConditions(command.QlWhere)
	objects := s.sqlObject(from)
	returns := fmt.Sprintf("%s AS result", objects)
	if len(command.Returns) > 0 {
		returns = ""
		for _, fld := range command.Returns {
			returns = strs.Append(returns, fld.Name, ", ")
		}
	}
	result := "DELETE FROM %s\nWHERE %s\nRETURNING\n%s;"
	return fmt.Sprintf(result, from.Model.Table, where, returns)
}

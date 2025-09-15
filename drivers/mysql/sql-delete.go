package mysql

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
func (s *Mysql) sqlDelete(command *jdb.Command) string {
	from := command.From
	where := whereConditions(command.QlWhere)
	objects := s.sqlObject(from.GetFrom())
	returns := fmt.Sprintf("%s AS result", objects)
	if len(command.Returns) > 0 {
		returns = ""
		for _, fld := range command.Returns {
			returns = strs.Append(returns, fld.Name, ", ")
		}
	}
	result := "DELETE FROM %s\nWHERE %s\nRETURNING\n%s;"
	return fmt.Sprintf(result, tableName(from), where, returns)
}

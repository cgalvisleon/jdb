package mysql

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* sqlFrom
* @param froms *jdb.QlFroms
* @return string
**/
func (s *Mysql) sqlFrom(froms *jdb.QlFroms) string {
	if len(froms.Froms) == 0 {
		return ""
	}

	from := froms.Froms[0]
	def := s.tableAs(from)
	result := strs.Format("FROM %s", def)

	return result
}

/**
* tableAs
* @param from *jdb.QlFrom
* @return string
**/
func (s *Mysql) tableAs(from *jdb.QlFrom) string {
	if from == nil {
		return ""
	}

	table := tableName(from.Model)
	return strs.Append(table, from.As, " AS ")
}

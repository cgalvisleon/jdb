package sqlite

import (
	"fmt"

	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* sqlWhere
* @param where *jdb.QlWhere
* @return string
**/
func (s *SqlLite) sqlWhere(where *jdb.QlWhere) string {
	if where == nil {
		return ""
	}

	if len(where.Wheres) == 0 {
		return ""
	}

	result := whereConditions(where)
	result = strs.Append("WHERE", result, " ")

	return result
}

/**
* whereConditions
* @param where *jdb.QlWhere
* @return string
**/
func whereConditions(where *jdb.QlWhere) string {
	result := ""
	for _, con := range where.Wheres {
		def := whereCondition(con)
		conector := whereConnector(con.Connector)
		result = strs.Append(result, def, conector)
	}

	return result
}

/**
* whereCondition
* @param con *jdb.QlCondition
* @return string
**/
func whereCondition(con *jdb.QlCondition) string {
	if con == nil {
		return ""
	}

	key := whereValue(con.Field)
	values := whereValue(con.Value)
	def := whereOperator(con, values)
	return fmt.Sprintf("%v%v", key, def)
}

/**
* whereValue
* @param val interface{}
* @return string
**/
func whereValue(val interface{}) string {
	switch v := val.(type) {
	case jdb.Field:
		return asField(v)
	case *jdb.Field:
		return asField(*v)
	case []interface{}:
		var result string
		for _, w := range v {
			val := whereValue(w)
			result = strs.Append(result, fmt.Sprintf(`%v`, val), ",")
		}
		return result
	default:
		return fmt.Sprintf(`%v`, jdb.Quote(v))
	}
}

/**
* whereOperator
* @param condition *jdb.QlCondition
* @param val interface{}
* @return string
**/
func whereOperator(condition *jdb.QlCondition, val interface{}) string {
	switch condition.Operator {
	case jdb.Equal:
		return fmt.Sprintf("=%v", val)
	case jdb.Neg:
		return fmt.Sprintf("!=%v", val)
	case jdb.In:
		return fmt.Sprintf(" IN (%v)", val)
	case jdb.Like:
		return fmt.Sprintf(" ILIKE %v", val)
	case jdb.More:
		return fmt.Sprintf(">%v", val)
	case jdb.Less:
		return fmt.Sprintf("<%v", val)
	case jdb.MoreEq:
		return fmt.Sprintf(">=%v", val)
	case jdb.LessEq:
		return fmt.Sprintf("<=%v", val)
	case jdb.Between:
		return fmt.Sprintf(" BETWEEN (%v)", val)
	case jdb.IsNull:
		return " IS NULL"
	case jdb.NotNull:
		return " IS NOT NULL"
	case jdb.Search:
		return fmt.Sprintf(" @@ to_tsquery('%s', %v)", condition.Language, val)
	default:
		return ""
	}
}

func whereConnector(con jdb.Connector) string {
	switch con {
	case jdb.And:
		return "\nAND "
	case jdb.Or:
		return "\nOR "
	default:
		return ""
	}
}

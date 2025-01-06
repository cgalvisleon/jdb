package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func whereOperator(op jdb.Operator, val interface{}) string {
	switch op {
	case jdb.Equal:
		return strs.Format("=%v", val)
	case jdb.Neg:
		return strs.Format("!=%v", val)
	case jdb.In:
		return strs.Format(" IN (%v)", val)
	case jdb.Like:
		return strs.Format(" ILIKE %v", val)
	case jdb.More:
		return strs.Format(">%v", val)
	case jdb.Less:
		return strs.Format("<%v", val)
	case jdb.MoreEq:
		return strs.Format(">=%v", val)
	case jdb.LessEq:
		return strs.Format("<=%v", val)
	case jdb.Between:
		return strs.Format(" BETWEEN (%v)", val)
	case jdb.IsNull:
		return " IS NULL"
	case jdb.NotNull:
		return " IS NOT NULL"
	case jdb.Search:
		if val == "%"+"%" {
			val = "%"
		}
		return strs.Format(" ILIKE %v", val)
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

func whereFilter(where *jdb.LinqWhere) string {
	if where == nil {
		return ""
	}

	key := where.GetKey()
	values := where.GetValue(where.Value)
	return strs.Format("%v%v", key, whereOperator(where.Operator, values))
}

func whereFilters(wheres []*jdb.LinqWhere) string {
	result := ""
	for _, w := range wheres {
		def := whereFilter(w)
		result = strs.Append(result, def, whereConnector(w.Conector))
	}

	return result
}

func (s *Postgres) queryWhere(wheres []*jdb.LinqWhere) string {
	result := whereFilters(wheres)
	result = strs.Append("WHERE", result, " ")

	return result
}

package postgres

import (
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func whereOperator(where *jdb.QlWhere, val interface{}) string {
	switch where.Operator {
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
		return strs.Format(" @@ to_tsquery('%s', %v)", where.Language, val)
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

func whereValue(val interface{}) string {
	adField := func(f *jdb.Field) string {
		switch f.Column.TypeColumn {
		case jdb.TpColumn:
			def := strs.Append(f.As, f.Field, ".")
			return strs.Format(`%s`, def)
		case jdb.TpAtribute:
			def := strs.Append(f.As, f.Field, ".")
			return strs.Format(`%s#>>'{%s}'`, def, f.Name)
		default:
			return ""
		}
	}

	switch v := val.(type) {
	case *jdb.QlSelect:
		return adField(v.Field)
	case *jdb.Field:
		return adField(v)
	case []interface{}:
		var result string
		for _, w := range v {
			val := whereValue(w)
			result = strs.Append(result, strs.Format(`%v`, val), ",")
		}
		return result
	default:
		return strs.Format(`%v`, utility.Quote(v))
	}
}

func whereKey(val interface{}) string {
	return whereValue(val)
}

func whereFilter(where *jdb.QlWhere) string {
	if where == nil {
		return ""
	}

	key := whereKey(where.Key)
	values := whereValue(where.Values)
	def := whereOperator(where, values)
	return strs.Format("%v%v", key, def)
}

func whereFilters(wheres []*jdb.QlWhere) string {
	result := ""
	for _, w := range wheres {
		def := whereFilter(w)
		conector := whereConnector(w.Conector)
		result = strs.Append(result, def, conector)
	}

	return result
}

func (s *Postgres) sqlWhere(wheres []*jdb.QlWhere) string {
	result := whereFilters(wheres)
	result = strs.Append("WHERE", result, " ")

	return result
}

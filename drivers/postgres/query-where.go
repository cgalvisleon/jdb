package postgres

import (
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func whereVal(val interface{}) interface{} {
	switch v := val.(type) {
	case *jdb.LinqSelect:
		return colName(v)
	default:
		return utility.Quote(v)
	}
}

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
		return " AND "
	case jdb.Or:
		return " OR "
	default:
		return ""
	}
}

func (s *Postgres) queryWhere(linq *jdb.Linq) string {
	result := ""
	for i, w := range linq.Wheres {
		a := whereVal(w.A)
		b := whereVal(w.B)
		def := strs.Format("%s%s", a, whereOperator(w.Operator, b))
		if i == 0 {
			result = def
			continue
		}

		result = strs.Append(result, def, whereConnector(w.Conector))
	}

	return result
}

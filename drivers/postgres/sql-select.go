package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func selectField(field *jdb.Field) string {
	agregaction := func(val string) string {
		switch field.Agregation {
		case jdb.AgregationSum:
			val = strs.Format(`SUM(%s)`, val)
		case jdb.AgregationCount:
			val = strs.Format(`COUNT(%s)`, val)
		case jdb.AgregationAvg:
			val = strs.Format(`AVG(%s)`, val)
		case jdb.AgregationMin:
			val = strs.Format(`MIN(%s)`, val)
		case jdb.AgregationMax:
			val = strs.Format(`MAX(%s)`, val)
		}

		return val
	}

	result := strs.Append("", field.As, "")
	switch field.Column.TypeColumn {
	case jdb.TpColumn:
		result = strs.Append(result, field.Field, ".")
		result = agregaction(result)
	case jdb.TpAtribute:
		result = strs.Append(result, field.Field, ".")
		result = strs.Format(`%s#>>'{%s}'`, result, field.Name)
		result = strs.Format(`COALESCE(%s, %v)`, result, field.Column.DefaultQuote())
		result = agregaction(result)
	}

	return result
}

func jsonbBuildObject(result, obj string) string {
	if len(obj) == 0 {
		return result
	}

	return strs.Append(result, strs.Format("jsonb_build_object(\n%s)", obj), "||\n")
}

func (s *Postgres) sqlData(frm *jdb.QlFrom, selects []*jdb.QlSelect) string {
	result := ""
	if frm != nil && frm.SourceField != nil {
		def := strs.Append("", frm.As, ".")
		def = strs.Append(def, frm.SourceField.Name, ".")
		def = strs.Format(`%s`, def)
		result = strs.Append(result, def, "")
	}
	l := 20
	n := 0
	obj := ""
	for _, sel := range selects {
		n++
		if sel.Field.Column == sel.From.SourceField {
			continue
		} else if sel.TypeColumn() == jdb.TpColumn {
			def := selectField(sel.Field)
			def = strs.Format(`'%s', %s`, sel.Field.Alias, def)
			obj = strs.Append(obj, def, ",\n")
		} else if sel.TypeColumn() == jdb.TpAtribute {
			def := selectField(sel.Field)
			def = strs.Format(`'%s', %s`, sel.Field.Alias, def)
			obj = strs.Append(obj, def, ",\n")
		}
		if n == l {
			result = jsonbBuildObject(result, obj)
			obj = ""
			n = 0
		}
	}
	if n > 0 {
		result = jsonbBuildObject(result, obj)
	}

	return result
}

func (s *Postgres) sqlDataOrders(frm *jdb.QlFrom, selects []*jdb.QlSelect, as string, orders []*jdb.QlOrder) string {
	result := s.sqlData(frm, selects)
	result = strs.Append(result, as, " AS ")
	for _, ord := range orders {
		def := selectField(ord.Field)
		result = strs.Append(result, def, ",\n")
	}

	return result
}

func (s *Postgres) sqlColumns(selects []*jdb.QlSelect) string {
	result := ""
	for _, sel := range selects {
		if sel.TypeColumn() == jdb.TpColumn {
			def := selectField(sel.Field)
			if sel.Field.Agregation != jdb.Nag {
				def = strs.Format(`%s AS %s`, def, sel.Field.Alias)
			}
			result = strs.Append(result, def, ",\n")
		} else if sel.TypeColumn() == jdb.TpAtribute {
			def := selectField(sel.Field)
			def = strs.Format(`%s AS %s`, def, sel.Field.Alias)
			result = strs.Append(result, def, ",\n")
		}
	}

	return result
}

func (s *Postgres) sqlSelect(ql *jdb.Ql) string {
	froms := ql.Froms
	if len(froms) == 0 {
		return ""
	}

	var result string
	if ql.TypeSelect == jdb.Data {
		result = s.sqlDataOrders(nil, ql.Selects, string(jdb.SourceField), ql.Orders)
	} else {
		result = s.sqlColumns(ql.Selects)
	}

	result = strs.Append("\nSELECT DISTINCT", result, "\n")

	return result
}

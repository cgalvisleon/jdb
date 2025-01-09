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
		result = strs.Append(result, field.Name, ".")
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

func (s *Postgres) sqlData(selects []*jdb.LinqSelect, orders []*jdb.LinqOrder) string {
	result := ""
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

	result = strs.Append(result, jdb.SourceField.Up(), " AS ")
	for _, ord := range orders {
		def := selectField(ord.Field)
		result = strs.Append(result, def, ",\n")
	}

	return result
}

func (s *Postgres) sqlColumns(tp jdb.TypeSelect, selects []*jdb.LinqSelect, orders []*jdb.LinqOrder) string {
	if tp == jdb.Data {
		return s.sqlData(selects, orders)
	}

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

func (s *Postgres) sqlSelect(linq *jdb.Linq) string {
	froms := linq.Froms
	if len(froms) == 0 {
		return ""
	}

	var selects = []*jdb.LinqSelect{}
	for _, frm := range froms {
		selects = append(selects, frm.Selects...)
	}

	if len(selects) == 0 {
		frm := froms[0]
		for _, col := range frm.Columns {
			field := frm.GetField(col.Name)
			if field != nil {
				selects = append(selects, &jdb.LinqSelect{
					From:  frm,
					Field: field,
				})
			}
		}
		selects = append(selects, frm.Selects...)
	}

	result := s.sqlColumns(linq.TypeSelect, selects, linq.Orders)
	result = strs.Append("\nSELECT DISTINCT", result, "\n")

	return result
}

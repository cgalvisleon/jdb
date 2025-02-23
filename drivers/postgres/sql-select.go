package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlSelect(ql *jdb.Ql) string {
	if len(ql.Froms.Froms) == 0 {
		return ""
	}

	var result string
	if ql.TypeSelect == jdb.Select {
		result = s.sqlColumns(ql.Selects)
	} else {
		result = s.sqlObjectOrders(ql.Selects, ql.Orders)
	}

	result = strs.Append("\nSELECT DISTINCT", result, "\n")

	return result
}

func asField(field jdb.Field) string {
	setAgregaction := func(val string) string {
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
		result = setAgregaction(result)
	case jdb.TpAtribute:
		result = strs.Append(result, field.Source, ".")
		result = strs.Format(`%s#>>'{%s}'`, result, field.Name)
		result = strs.Format(`COALESCE(%s, %v)`, result, field.Column.DefaultQuote())
		result = setAgregaction(result)
	}

	return result
}

func selectAsField(field jdb.Field) string {
	result := asField(field)
	if field.Name == field.Alias {
		return result
	}

	return strs.Append(result, field.Alias, " AS ")
}

func jsonbBuildObject(result, obj string) string {
	if len(obj) == 0 {
		return result
	}

	return strs.Append(result, strs.Format("jsonb_build_object(\n%s)", obj), "||\n")
}

func (s *Postgres) sqlObject(selects []*jdb.Field) string {
	result := ""
	l := 20
	if s.version >= 13 {
		l = 100
	}
	n := 0
	obj := ""
	sourceField := make([]*jdb.Field, 0)
	for _, fld := range selects {
		n++
		if fld.Hidden {
			continue
		}
		col := fld.Column
		if col == col.Model.SourceField {
			sourceField = append(sourceField, fld)
			continue
		}
		def := selectAsField(*fld)
		def = strs.Format(`'%s', %s`, fld.Alias, def)
		obj = strs.Append(obj, def, ",\n")

		if n == l {
			result = jsonbBuildObject(result, obj)
			obj = ""
			n = 0
		}
	}
	if n > 0 {
		result = jsonbBuildObject(result, obj)
	}
	sources := ""
	for i := 0; i < len(sourceField); i++ {
		fld := sourceField[i]
		def := selectAsField(*fld)
		sources = strs.Append(sources, def, "||\n")
		if i == len(sourceField)-1 {
			result = strs.Format(`%s||%s`, def, result)
		}
	}

	return result
}

func (s *Postgres) sqlObjectOrders(selects []*jdb.Field, orders *jdb.QlOrder) string {
	result := s.sqlObject(selects)
	result = strs.Append(result, "result", " AS ")
	for _, ord := range orders.Asc {
		def := selectAsField(*ord)
		result = strs.Append(result, def, ",\n")
	}
	for _, ord := range orders.Desc {
		def := selectAsField(*ord)
		result = strs.Append(result, def, ",\n")
	}

	return result
}

func (s *Postgres) sqlColumns(selects []*jdb.Field) string {
	result := ""
	for _, fld := range selects {
		if fld.Hidden {
			continue
		}
		def := selectAsField(*fld)
		result = strs.Append(result, def, ",\n")
	}

	return result
}

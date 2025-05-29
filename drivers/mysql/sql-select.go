package mysql

import (
	"slices"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* sqlSelect
* @param ql *jdb.Ql
* @return string
**/
func (s *Mysql) sqlSelect(ql *jdb.Ql) string {
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

/**
* asField
* @param field jdb.Field
* @return string
**/
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

/**
* aliasAsField
* @param field jdb.Field
* @return string
**/
func aliasAsField(field jdb.Field) string {
	result := asField(field)
	if field.Name == field.Alias {
		return result
	}

	return strs.Append(result, field.Alias, " AS ")
}

/**
* jsonBuildObject
* @param result, obj string
* @return string
**/
func jsonBuildObject(result, obj string) string {
	if len(obj) == 0 {
		return result
	}

	return strs.Append(result, strs.Format("jsonb_build_object(\n%s)", obj), "||\n")
}

/**
* sqlObject
* @param from *jdb.QlFrom
* @return string
**/
func (s *Mysql) sqlObject(from *jdb.QlFrom) string {
	var selects = []*jdb.Field{}
	for _, col := range from.Columns {
		field := col.GetField()
		if field == nil {
			continue
		}
		if field.Column == nil {
			continue
		}
		if slices.Contains([]jdb.TypeColumn{jdb.TpColumn}, field.Column.TypeColumn) {
			field.As = from.As
			selects = append(selects, field)
		}
	}

	return s.sqlBuildObject(selects)
}

/**
* sqlBuildObject
* @param selects []*jdb.Field
* @return string
**/
func (s *Mysql) sqlBuildObject(selects []*jdb.Field) string {
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
		def := asField(*fld)
		def = strs.Format(`'%s', %s`, fld.Alias, def)
		obj = strs.Append(obj, def, ",\n")

		if n == l {
			result = jsonBuildObject(result, obj)
			obj = ""
			n = 0
		}
	}
	if n > 0 {
		result = jsonBuildObject(result, obj)
	}
	sources := ""
	for i := 0; i < len(sourceField); i++ {
		fld := sourceField[i]
		def := aliasAsField(*fld)
		sources = strs.Append(sources, def, "||\n")
		if i == len(sourceField)-1 {
			result = strs.Format(`%s||%s`, def, result)
		}
	}

	return result
}

/**
* sqlObjectOrders
* @param selects []*jdb.Field, orders *jdb.QlOrder
* @return string
**/
func (s *Mysql) sqlObjectOrders(selects []*jdb.Field, orders *jdb.QlOrder) string {
	result := s.sqlBuildObject(selects)
	result = strs.Append(result, "result", " AS ")
	for _, ord := range orders.Asc {
		def := aliasAsField(*ord)
		result = strs.Append(result, def, ",\n")
	}
	for _, ord := range orders.Desc {
		def := aliasAsField(*ord)
		result = strs.Append(result, def, ",\n")
	}

	return result
}

/**
* sqlColumns
* @param selects []*jdb.Field
* @return string
**/
func (s *Mysql) sqlColumns(selects []*jdb.Field) string {
	result := ""
	for _, fld := range selects {
		if fld.Hidden {
			continue
		}
		def := aliasAsField(*fld)
		result = strs.Append(result, def, ",\n")
	}

	return result
}

/**
* Select
* @param ql *jdb.Ql
* @return et.Items, error
**/
func (s *Mysql) Select(ql *jdb.Ql) (et.Items, error) {
	ql.Sql = ""
	ql.Sql = strs.Append(ql.Sql, s.sqlSelect(ql), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlFrom(ql.Froms), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlJoin(ql.Joins), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlWhere(ql.QlWhere), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlGroupBy(ql), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlHaving(ql), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlOrderBy(ql), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlLimit(ql), "\n")
	ql.Sql = strs.Format(`%s;`, ql.Sql)

	if ql.IsDebug {
		console.Debug(ql.Sql)
	}

	if ql.TypeSelect == jdb.Source {
		result, err := jdb.DataTx(ql.Tx(), s.db, "result", ql.Sql)
		if err != nil {
			return et.Items{}, err
		}

		return result, nil
	}

	result, err := jdb.QueryTx(ql.Tx(), s.db, ql.Sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func colName(sl *jdb.LinqSelect) string {
	result := strs.Append("", sl.From.As, "")
	switch sl.TypeColumn() {
	case jdb.TpColumn:
		result = strs.Append(result, sl.Field.Name, ".")
	case jdb.TpAtribute:
		result = strs.Append(result, sl.Field.Name, ".")
		result = strs.Format(`%s#>>'{%s}'`, result, sl.Field.Atrib)
	}

	return result
}

func jsonbBuildObject(result, obj string) string {
	return strs.Append(result, strs.Format("jsonb_build_object(\n%s)", obj), "||\n")
}

func (s *Postgres) queryColumns(froms []*jdb.LinqFrom) string {
	if len(froms) == 0 {
		return ""
	}

	var result string
	for _, frm := range froms {
		l := 20
		n := 0
		obj := ""
		for _, sel := range frm.Selects {
			n++
			if sel.TypeColumn() == jdb.TpColumn && sel.Field.Column != sel.From.SourceField {
				def := colName(sel)
				obj = strs.Append(obj, strs.Format(`'%s', %s`, sel.Field.Name, def), ",\n")
			} else if sel.TypeColumn() == jdb.TpAtribute {
				def := colName(sel)
				def = strs.Format(`COALESCE(%s, %v)`, def, sel.Field.Column.DefaultQuote())
				def = strs.Format(`'%s', %s`, sel.Field.Atrib, def)
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
	}

	if len(result) == 0 {
		frm := froms[0]
		if frm.SourceField != nil {
			result = strs.Format("%s.%s", frm.As, frm.SourceField.Up())
		}
		l := 20
		n := 0
		obj := ""
		for _, col := range frm.Columns {
			n++
			if col.TypeColumn == jdb.TpColumn && col != col.Model.SourceField {
				obj = strs.Append(obj, strs.Format(`'%s', %s.%s`, col.Low(), frm.As, col.Up()), ",\n")
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
	}

	return result
}

func (s *Postgres) querySelect(froms []*jdb.LinqFrom) string {
	columns := s.queryColumns(froms)
	if len(columns) == 0 {
		return ""
	}

	result := "\nSELECT DISTINCT"
	result = strs.Append(result, columns, "\n")
	result = strs.Append(result, strs.Format(" AS %s", jdb.SourceField.Up()), "")

	return result
}

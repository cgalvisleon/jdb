package postgres

import (
	"fmt"
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

/**
* Query
* @param query *jdb.Ql
* @return (string, error)
**/
func (s *Driver) buildQuery(ql *jdb.Ql) (string, error) {
	query := ql.ToJson()

	if ql.IsDebug {
		logs.Debug("query:", query.ToString())
	}

	sql, err := s.buildSelect(query)
	if err != nil {
		return "", err
	}

	sql = fmt.Sprintf("SELECT %s", sql)
	def, err := s.buildFrom(query)
	if err != nil {
		return "", err
	}

	def = fmt.Sprintf("FROM %s", def)
	sql = strs.Append(sql, def, "\n")
	def, err = s.buildJoins(query)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("JOIN %s", def)
		sql = strs.Append(sql, def, "\n")
	}

	where := query.ArrayJson("where")
	if len(where) > 0 {
		def, err = s.buildWhere(where)
		if err != nil {
			return "", err
		}

		if def != "" {
			def = fmt.Sprintf("WHERE %s", def)
			sql = strs.Append(sql, def, "\n")
		}
	}

	def, err = s.buildGroupBy(query)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("GROUP BY %s", def)
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildHaving(query)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("HAVING %s", def)
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildOrderBy(query)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("ORDER BY %s", def)
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildLimit(query)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	if query.Bool("exists") {
		return fmt.Sprintf("SELECT EXISTS(%s);", sql), nil
	} else {
		return fmt.Sprintf("%s;", sql), nil
	}
}

/**
* buildSelect
* @param query et.Json
* @return (string, error)
**/
func (s *Driver) buildSelect(query et.Json) (string, error) {
	result := ""

	if query.Bool("exists") {
		return "", nil
	}

	if query.Bool("count") {
		return "COUNT(*) AS all", nil
	}

	isDataSource := query.Bool("is_data_source")
	if isDataSource {
		atribs := query.Json("atribs")
		if atribs.IsEmpty() {
			result = fmt.Sprintf("\n%s", jdb.SOURCE)
		} else {
			for k, v := range atribs {
				def := fmt.Sprintf("\n'%s', %s", k, v)
				result = strs.Append(result, def, ", ")
			}

			if result != "" {
				result = fmt.Sprintf("\n\tjsonb_build_object(%s\n)", result)
			}
		}

		selects := query.Json("selects")
		if selects.IsEmpty() {
			hidden := query.ArrayStr("hidden")
			hidden = append(hidden, jdb.SOURCE)
			def := fmt.Sprintf("to_jsonb(A) - ARRAY[%s]", strings.Join(hidden, ", "))
			result = strs.Append(result, def, "||")
		} else {
			sel := ""
			for k := range selects {
				v := selects.String(k)
				def := fmt.Sprintf("\n'%s',  %s", k, v)
				if v == "" {
					def = fmt.Sprintf("\n'%s',  %s", k, k)
				}
				sel = strs.Append(sel, def, ", ")
			}

			if sel != "" {
				result = fmt.Sprintf("%s||jsonb_build_object(%s\n)", result, sel)
			}
		}

		return fmt.Sprintf("%s AS result", result), nil
	}

	selects := query.Json("selects")
	if selects.IsEmpty() {
		hidden := query.ArrayStr("hidden")
		if len(hidden) > 0 {
			result += fmt.Sprintf("to_jsonb(A) - ARRAY[%s]", strings.Join(hidden, ", "))
		} else {
			result += "A.*"
		}
	} else {
		for k := range selects {
			v := selects.String(k)
			def := fmt.Sprintf("\n%s AS %s", v, k)
			if k == v {
				def = fmt.Sprintf("\n%s", v)
			} else if v == "" {
				def = fmt.Sprintf("\n%s", v)
			}

			result = strs.Append(result, def, ", ")
		}
	}

	return result, nil
}

/**
* buildFrom
* @param query et.Json
* @return (string, error)
**/
func (s *Driver) buildFrom(query et.Json) (string, error) {
	result := ""

	froms := query.ArrayJson("from")
	if len(froms) == 0 {
		return result, fmt.Errorf(jdb.MSG_FROM_REQUIRED)
	}

	for _, v := range froms {
		as := v.Str("as")
		table := v.Str("table")
		def := fmt.Sprintf("%s AS %s", table, as)
		if as == table {
			def = fmt.Sprintf("%s", table)
		}

		result = strs.Append(result, def, ", ")
		break
	}

	return result, nil
}

/**
* buildJoins
* @param query et.Json
* @return (string, error)
**/
func (s *Driver) buildJoins(query et.Json) (string, error) {
	result := ""

	joins := query.ArrayJson("joins")
	if len(joins) == 0 {
		return result, nil
	}

	for _, v := range joins {
		def, err := s.buildFrom(v)
		if err != nil {
			return "", err
		}

		result = strs.Append(result, def, "")
		on := v.ArrayJson("on")
		def, err = s.buildWhere(on)
		if err != nil {
			return "", err
		}

		if def != "" {
			def = fmt.Sprintf("ON %s", def)
		}

		result = strs.Append(result, def, " ")
	}

	return fmt.Sprintf("%s", result), nil
}

/**
* buildWhere
* @param wheres []et.Json
* @return (string, error)
**/
func (s *Driver) buildWhere(wheres []et.Json) (string, error) {
	result := ""

	condition := func(condition et.Json) string {
		for k, v := range condition {
			switch k {
			case "eq":
				def := fmt.Sprintf("= %v", v)
				return def
			case "neg":
				def := fmt.Sprintf("!= %v", v)
				return def
			case "less":
				def := fmt.Sprintf("< %v", v)
				return def
			case "less_eq":
				def := fmt.Sprintf("<= %v", v)
				return def
			case "more":
				def := fmt.Sprintf("> %v", v)
				return def
			case "more_eq":
				def := fmt.Sprintf(">= %v", v)
				return def
			case "like":
				def := fmt.Sprintf("LIKE %v", v)
				return def
			case "ilike":
				def := fmt.Sprintf("ILIKE %v", v)
				return def
			case "in":
				def := fmt.Sprintf("IN %v", v)
				return def
			case "not_in":
				def := fmt.Sprintf("NOT IN %v", v)
				return def
			case "is":
				def := fmt.Sprintf("IS %v", v)
				return def
			case "is_not":
				def := fmt.Sprintf("IS NOT %v", v)
				return def
			case "null":
				def := fmt.Sprintf("IS NULL")
				return def
			case "not_null":
				def := fmt.Sprintf("IS NOT NULL")
				return def
			case "between":
				vals := condition.Array(k)
				def := fmt.Sprintf("BETWEEN %v AND %v", vals[0], vals[1])
				return def
			case "not_between":
				vals := condition.Array(k)
				def := fmt.Sprintf("NOT BETWEEN %v AND %v", vals[0], vals[1])
				return def
			}
		}

		return ""
	}

	append := func(cond, connect string) string {
		if result == "" {
			return cond
		} else {
			connect = fmt.Sprintf("\n%s ", connect)
			return strs.Append(result, cond, connect)
		}
	}

	logs.Log("WHERE clauses:", et.ArrayJsonToString(wheres))

	for _, w := range wheres {
		for k := range w {
			if map[string]bool{"and": true, "or": true}[strs.Lowcase(k)] {
				andOr := w.ArrayJson(k)
				for _, field := range andOr {
					for name := range field {
						v := field.Json(name)
						cond := condition(v)
						def := fmt.Sprintf("%s %s", name, cond)
						result = append(def, strs.Uppcase(k))
					}
				}
			} else {
				v := w.Json(k)
				cond := condition(v)
				def := fmt.Sprintf("%s %s", k, cond)
				result = append(def, "AND")
			}
		}
	}

	return result, nil
}

/**
* buildGroupBy
* @param query et.Json
* @return (string, error)
**/
func (s *Driver) buildGroupBy(query et.Json) (string, error) {
	result := ""

	groupBy := query.ArrayStr("group_by")
	if len(groupBy) == 0 {
		return result, nil
	}

	for _, v := range groupBy {
		def := fmt.Sprintf("%s", v)
		result = strs.Append(result, def, ", ")
	}

	return result, nil
}

/**
* buildHaving
* @param query et.Json
* @return (string, error)
**/
func (s *Driver) buildHaving(query et.Json) (string, error) {
	result := ""

	having := query.ArrayJson("having")
	if len(having) == 0 {
		return result, nil
	}

	result, err := s.buildWhere(having)
	if err != nil {
		return "", err
	}

	return result, nil
}

/**
* buildOrderBy
* @param query et.Json
* @return (string, error)
**/
func (s *Driver) buildOrderBy(query et.Json) (string, error) {
	result := ""

	orderBy := query.Json("order_by")
	if orderBy.IsEmpty() {
		return result, nil
	}

	asc := orderBy.ArrayStr("asc")
	desc := orderBy.ArrayStr("desc")
	ascs := ""
	for _, v := range asc {
		ascs = strs.Append(ascs, v, ", ")
	}

	if ascs != "" {
		ascs = fmt.Sprintf(`%s ASC`, ascs)
	}

	descs := ""
	for _, v := range desc {
		descs = strs.Append(descs, v, ", ")
	}

	if descs != "" {
		descs = fmt.Sprintf(`%s DESC`, descs)
	}

	result = strs.Append(ascs, descs, ", ")
	return result, nil
}

/**
* buildLimit
* @param query et.Json
* @return (string, error)
**/
func (s *Driver) buildLimit(query et.Json) (string, error) {
	result := ""

	limit := query.Json("limit")
	if limit.IsEmpty() {
		return result, nil
	}

	limitRows := query.ValInt(1000, "max_rows")
	page := limit.Int("page")
	rows := limit.ValInt(limitRows, "rows")

	if page == 0 {
		result = fmt.Sprintf("LIMIT %d", rows)
		return result, nil
	}

	offset := (page - 1) * rows
	result = fmt.Sprintf("%d OFFSET %d", rows, offset)
	return result, nil
}

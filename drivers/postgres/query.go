package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) buildSelect(query et.Json) (string, error) {
	result := ""
	tp := query.String("type")
	if tp == jdb.TpObject {
		atribs := query.Json("atribs")
		if atribs.IsEmpty() {
			result = fmt.Sprintf("\n\t%s", jdb.SOURCE)
		} else {
			for k, v := range atribs {
				def := fmt.Sprintf("\n\t'%s', %s", v, k)
				result = strs.Append(result, def, ", ")
			}

			if result != "" {
				result = fmt.Sprintf("\n\tjsonb_build_object(%s\n\t)", result)
			}
		}

		selects := query.Json("selects")
		if selects.IsEmpty() {
			def := fmt.Sprintf("to_jsonb(A) - '%s'", jdb.SOURCE)
			result = strs.Append(result, def, "||")
		} else {
			sel := ""
			for k := range selects {
				v := selects.String(k)
				def := fmt.Sprintf("\n\t'%s',  %s", v, k)
				sel = strs.Append(sel, def, ", ")
			}

			if sel != "" {
				result = fmt.Sprintf("%s||jsonb_build_object(%s\n\t)", result, sel)
			}
		}

		return fmt.Sprintf("%s AS result", result), nil
	}

	selects := query.Json("selects")
	if selects.IsEmpty() {
		result += "*"
	} else {
		for k := range selects {
			v := selects.String(k)
			def := fmt.Sprintf("\n\t%s AS %s", k, v)
			if k == v {
				def = fmt.Sprintf("\n\t%s", k)
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
func (s *Postgres) buildFrom(query et.Json) (string, error) {
	froms := query.Json("from")
	if froms.IsEmpty() {
		return "", fmt.Errorf(jdb.MSG_FROM_REQUIRED)
	}

	result := ""
	for k := range froms {
		v := froms.String(k)
		def := fmt.Sprintf("%s AS %s", k, v)
		if k == v {
			def = fmt.Sprintf("%s", k)
		}

		result = strs.Append(result, def, ", ")
	}

	return result, nil
}

/**
* buildJoins
* @param query et.Json
* @return (string, error)
**/
func (s *Postgres) buildJoins(query et.Json) (string, error) {
	joins := query.ArrayJson("joins")
	if len(joins) == 0 {
		return "", nil
	}

	result := ""
	for _, v := range joins {
		def, err := s.buildFrom(v)
		if err != nil {
			return "", err
		}

		result = strs.Append(result, def, "")
		on := v.Json("on")
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
* @param wheres et.Json
* @return (string, error)
**/
func (s *Postgres) buildWhere(wheres et.Json) (string, error) {
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
			case "exists":
				def := fmt.Sprintf("EXISTS %v", v)
				return def
			case "not_exists":
				def := fmt.Sprintf("NOT EXISTS %v", v)
				return def
			}
		}

		return ""
	}

	result := ""
	append := func(cond, connect string) string {
		if result == "" {
			return cond
		} else {
			connect = fmt.Sprintf("\n\t%s ", connect)
			return strs.Append(result, cond, connect)
		}
	}

	for k := range wheres {
		if map[string]bool{"AND": true, "OR": true}[strs.Uppcase(k)] {
			andOr := wheres.Json(k)
			for f := range andOr {
				v := andOr.Json(f)
				cond := condition(v)
				def := fmt.Sprintf("%s %s", f, cond)
				result = append(def, strs.Uppcase(k))
			}
		} else {
			v := wheres.Json(k)
			cond := condition(v)
			def := fmt.Sprintf("%s %s", k, cond)
			result = append(def, "AND")
		}
	}

	return result, nil
}

/**
* buildGroupBy
* @param query et.Json
* @return (string, error)
**/
func (s *Postgres) buildGroupBy(query et.Json) (string, error) {
	groupBy := query.ArrayStr("group_by")
	if len(groupBy) == 0 {
		return "", nil
	}

	result := ""
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
func (s *Postgres) buildHaving(query et.Json) (string, error) {
	having := query.Json("having")
	if having.IsEmpty() {
		return "", nil
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
func (s *Postgres) buildOrderBy(query et.Json) (string, error) {
	orderBy := query.Json("order_by")
	if orderBy.IsEmpty() {
		return "", nil
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

	result := strs.Append(ascs, descs, ", ")
	return result, nil
}

/**
* buildLimit
* @param query et.Json
* @return (string, error)
**/
func (s *Postgres) buildLimit(query et.Json) (string, error) {
	limit := query.Json("limit")
	if limit.IsEmpty() {
		return "", nil
	}

	limitRows := envar.GetInt("LIMIT_ROWS", 100)
	page := limit.Int("page")
	rows := limit.ValInt(limitRows, "rows")

	if page == 0 {
		return fmt.Sprintf("%d", rows), nil
	}

	offset := (page - 1) * rows
	return fmt.Sprintf("%d OFFSET %d", rows, offset), nil
}

/**
* Query
* @param query *jdb.Ql
* @return (string, error)
**/
func (s *Postgres) buildQuery(query et.Json) (string, error) {
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
	sql = strs.Append(sql, def, "\n\t")
	def, err = s.buildJoins(query)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("JOIN %s", def)
		sql = strs.Append(sql, def, "\n\t")
	}

	where := query.Json("where")
	if !where.IsEmpty() {
		def, err = s.buildWhere(where)
		if err != nil {
			return "", err
		}

		if def != "" {
			def = fmt.Sprintf("WHERE %s", def)
			sql = strs.Append(sql, def, "\n\t")
		}
	}

	def, err = s.buildGroupBy(query)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("GROUP BY %s", def)
		sql = strs.Append(sql, def, "\n\t")
	}

	def, err = s.buildHaving(query)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("HAVING %s", def)
		sql = strs.Append(sql, def, "\n\t")
	}

	def, err = s.buildOrderBy(query)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("ORDER BY %s", def)
		sql = strs.Append(sql, def, "\n\t")
	}

	def, err = s.buildLimit(query)
	if err != nil {
		return "", err
	}

	if def != "" {
		def = fmt.Sprintf("LIMIT %s", def)
		sql = strs.Append(sql, def, "\n\t")
	}

	return sql, nil
}

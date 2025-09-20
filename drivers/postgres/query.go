package postgres

import (
	"fmt"

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
			result += jdb.SOURCE
		} else {
			for k, v := range atribs {
				def := fmt.Sprintf("(\n\t'%s', %s", v, k)
				result = strs.Append(result, def, ", ")
			}
		}
		if result != "" {
			result = fmt.Sprintf("jsonb_build_object(%s\n\t)", result)
		}

		sel := ""
		selects := query.Json("selects")
		for k := range selects {
			v := selects.String(k)
			def := fmt.Sprintf("\n\t'%s',  %s", v, k)
			sel = strs.Append(sel, def, ", ")
		}

		if sel != "" {
			result = fmt.Sprintf("%s||jsonb_build_object(%s\n\t)", result, sel)
		}

		return fmt.Sprintf("SELECT %s", result), nil
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

	result = fmt.Sprintf("SELECT %s", result)
	return result, nil
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

	return sql, nil
}

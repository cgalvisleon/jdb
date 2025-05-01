package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
)

/**
* GetDetailsTx
* @param tx *Tx, data *et.Json
* @return et.Json
**/
func (s *Ql) GetDetailsTx(tx *Tx, data *et.Json) *et.Json {
	for _, field := range s.Details {
		col := field.Column
		if col == nil {
			continue
		}

		switch col.TypeColumn {
		case TpGenerated:
			if col.GeneratedFunction != nil {
				col.GeneratedFunction(col, data)
			}
		case TpRelatedTo:
			if col.Detail == nil {
				continue
			}
			with := col.Detail.With
			if with == nil {
				continue
			}

			where := col.Detail.Where(*data)
			if s.IsDebug {
				console.Debug(where.ToString())
			}
			ql := From(with).
				setJoins(field.Joins).
				setWheres(where).
				setWheres(field.Where).
				setSelects(field.Select).
				setGroupBy(field.GroupBy...).
				setHavings(field.Havings).
				setOrderBy(field.OrderBy).
				setDebug(s.IsDebug)

			if field.TpResult == TpResult {
				result, err := ql.AllTx(tx)
				if err != nil {
					continue
				}

				data.Set(col.Name, result.Result)
			} else {
				all, err := ql.
					CountedTx(tx)
				if err != nil {
					continue
				}

				result, err := ql.
					Page(field.Page).
					RowsTx(tx, field.Rows)
				if err != nil {
					continue
				}

				data.Set(col.Name, result.ToList(all, field.Page, field.Rows))
			}
		case TpRollup:
			if col.Rollup == nil {
				continue
			}
			source := col.Rollup.Source
			if source == nil {
				continue
			}

			where := col.Rollup.Where(*data)
			if s.IsDebug {
				console.Debug(where.ToString())
			}
			ql := From(source).
				setJoins(field.Joins).
				setWheres(where).
				setWheres(field.Where).
				setSelects(field.Select).
				setGroupBy(field.GroupBy...).
				setHavings(field.Havings).
				setOrderBy(field.OrderBy).
				setDebug(s.IsDebug)

			fields := col.Rollup.Fields

			switch len(fields) {
			case 0:
				continue
			case 1:
				field, ok := fields[0].(string)
				if !ok {
					continue
				}
				result, err := ql.
					Data(field).
					OneTx(tx)
				if err != nil {
					continue
				}

				data.Set(col.Name, result.Result[field])
			default:
				result, err := ql.
					Data(fields...).
					OneTx(tx)
				if err != nil {
					continue
				}

				data.Set(col.Name, result.Result)
			}
		}
	}

	return data
}

/**
* Detail
* @param name string, selects []interface{}, joins []et.Json, where et.Json, groups []string, havings et.Json, orderBy et.Json, page, rows int
* @return *Ql
**/
func (s *Ql) Detail(name string, selects []interface{}, joins []et.Json, where et.Json, groups []string, havings et.Json, orderBy et.Json, page, rows int, tp TypeResult) *Ql {
	field := s.getField(name, false)
	if field == nil || field.Column == nil || field.Column.Detail == nil || field.Column.Detail.With == nil {
		return s
	}

	idx := slices.IndexFunc(s.Details, func(e *Field) bool { return e.asField() == field.asField() })
	if idx == -1 {
		s.Details = append(s.Details, field)
	}

	field.Select = selects
	field.Joins = joins
	field.Where = where
	field.GroupBy = groups
	field.Havings = havings
	field.OrderBy = orderBy
	field.Page = page
	field.Rows = rows
	field.TpResult = tp

	return s
}

/**
* setDetail
* @param params et.Json
* @return *Ql
**/
func (s *Ql) setDetail(params et.Json) *Ql {
	for name := range params {
		val := params.Json(name)
		selects := val.Array("select")
		joins := val.ArrayJson("join")
		where := val.Json("where")
		groups := val.ArrayStr("group_by")
		havings := val.Json("having")
		orderBy := val.Json("order_by")
		page := val.Int("page")
		rows := val.Int("rows")
		list := val.Bool("list")
		tp := TpResult
		if list {
			tp = TpList
		}

		s.Detail(name, selects, joins, where, groups, havings, orderBy, page, rows, tp)
	}

	return s
}

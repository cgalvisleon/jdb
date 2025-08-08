package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
)

/**
* GetDetailsTx
* @param tx *Tx, data et.Json
* @return error
**/
func (s *Ql) GetDetailsTx(tx *Tx, data et.Json) {
	sort := []*Field{}
	sort1 := []*Field{}
	sort2 := []*Field{}
	sort3 := []*Field{}
	for _, field := range s.Details {
		col := field.Column
		if col == nil {
			continue
		}

		switch col.TypeColumn {
		case TpCalc:
			sort1 = append(sort1, field)
		case TpRelatedTo:
			sort1 = append(sort1, field)
		case TpRollup:
			sort1 = append(sort1, field)
		case TpConcurrent:
			sort3 = append(sort3, field)
		default:
			sort2 = append(sort2, field)
		}
	}

	sort = append(sort, sort1...)
	sort = append(sort, sort2...)

	for _, field := range sort {
		col := field.Column
		switch col.TypeColumn {
		case TpCalc:
			if col.CalcFunction != nil {
				col.CalcFunction(data)
			}
		case TpAtribute:
			if col.CalcFunction != nil {
				col.CalcFunction(data)
			}
		case TpColumn:
			if col.CalcFunction != nil {
				col.CalcFunction(data)
			}
		case TpRelatedTo:
			if col.CalcFunction != nil {
				col.CalcFunction(data)
				continue
			}

			if col.Detail == nil {
				continue
			}

			with := col.Detail.With
			if with == nil {
				continue
			}

			where := col.Detail.GetWhere(data)
			if s.IsDebug {
				console.Debug("GetDetailsTx:", where.ToString())
			}

			ql := From(with).
				SetJoins(field.Joins).
				SetWheres(where).
				SetWheres(field.Where).
				SetSelects(field.Select).
				SetGroupBy(field.GroupBy...).
				SetHavings(field.Havings).
				SetOrderBy(field.OrderBy).
				SetDebug(s.IsDebug).
				prepare()

			idx := slices.IndexFunc(ql.Details, func(e *Field) bool { return e.Column == field.Column })
			if idx != -1 {
				ql.Details = append(ql.Details[:idx], ql.Details[idx+1:]...)
			}

			if field.TpResult == TpResult {
				result, err := ql.AllTx(tx)
				if err != nil {
					continue
				}
				if !result.Ok {
					data.Set(col.Name, et.Json{})
				} else {
					data.Set(col.Name, result.Result)
				}
			} else {
				all, err := ql.
					CountedTx(tx)
				if err != nil {
					continue
				}

				if all == 0 {
					data.Set(col.Name, et.Json{})
				} else {
					result, err := ql.
						Page(field.Page).
						RowsTx(tx, field.Rows)
					if err != nil {
						continue
					}

					data.Set(col.Name, result.ToList(all, field.Page, field.Rows))
				}
			}
		case TpRollup:
			if col.CalcFunction != nil {
				col.CalcFunction(data)
				continue
			}

			if col.Rollup == nil {
				continue
			}

			rollup := col.Rollup
			with := rollup.With
			if with == nil {
				continue
			}

			where := rollup.Where(data)
			if s.IsDebug {
				console.Debug("GetDetailsTx:", where.ToString())
			}

			ql := From(with).
				SetJoins(field.Joins).
				SetWheres(where).
				SetWheres(field.Where).
				SetSelects(field.Select...).
				SetGroupBy(field.GroupBy...).
				SetHavings(field.Havings).
				SetOrderBy(field.OrderBy).
				SetDebug(s.IsDebug).
				prepare()

			idx := slices.IndexFunc(ql.Details, func(e *Field) bool { return e.Column == field.Column })
			if idx != -1 {
				ql.Details = append(ql.Details[:idx], ql.Details[idx+1:]...)
			}

			if len(ql.Selects) == 0 {
				continue
			}

			result, err := ql.
				OneTx(tx)
			if err != nil {
				continue
			}

			if rollup.Show == ShowObject {
				object := et.Json{}
				for fkn, pkn := range rollup.Fk {
					val := data.Str(fkn)
					object.Set(pkn, val)
				}
				for key, val := range result.Result {
					object.Set(key, val)
				}

				data.Set(col.Name, object)
			} else {
				for _, val := range result.Result {
					data.Set(col.Name, val)
				}
			}
		}
	}

	for _, field := range sort3 {
		col := field.Column
		if col == nil {
			continue
		}

		if col.CalcFunction == nil {
			continue
		}

		s.wg.Add(1)
		go func(data et.Json) {
			defer s.wg.Done()
			col.CalcFunction(data)
		}(data)
	}

	s.wg.Wait()
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
		console.Ping()
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

		field := s.getField(name)
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
	}

	return s
}

package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/et"
)

/**
* GetDetails
* @return et.Json
**/
func (s *Ql) GetDetails(data *et.Json) *et.Json {
	for _, field := range s.Details {
		col := field.Column
		if col == nil {
			continue
		}

		switch col.TypeColumn {
		case TpGenerated:
			if col.generatedFunction != nil {
				col.generatedFunction(col, data)
			}
		case TpRelatedTo:
			if col.Detail == nil {
				continue
			}
			if len(col.Detail.Fk) <= 0 {
				continue
			}

			n := 0
			with := col.Detail.With
			if with == nil {
				continue
			}

			ql := From(with)
			for fkn, pk := range col.Detail.Fk {
				key := (*data)[pk]
				if key == nil {
					continue
				}

				if n == 0 {
					ql.Where(fkn).Eq(key)
				} else {
					ql.And(fkn).Eq(key)
				}
				n++
			}

			if field.TpResult == TpResult {
				result, err := ql.All()
				if err != nil {
					continue
				}
				data.Set(col.Name, result.Result)
			} else {
				all, err := ql.
					Counted()
				if err != nil {
					continue
				}

				result, err := ql.
					Page(field.Page).
					Rows(field.Rows)
				if err != nil {
					continue
				}

				data.Set(col.Name, result.ToList(all, field.Page, field.Rows))
			}
		case TpRollup:
			if col.Rollup == nil {
				continue
			}
			if len(col.Rollup.Fk) <= 0 {
				continue
			}

			n := 0
			source := col.Rollup.Source
			if source == nil {
				continue
			}

			ql := From(source)
			for fkn, pk := range col.Rollup.Fk {
				key := (*data)[pk]
				if key == nil {
					continue
				}

				if n == 0 {
					ql.Where(fkn).Eq(key)
				} else {
					ql.And(fkn).Eq(key)
				}
				n++
			}

			props := make([]string, 0)
			props = append(props, col.Rollup.Props...)

			switch len(props) {
			case 0:
				continue
			case 1:
				prop := props[0]
				result, err := ql.
					Data(prop).
					One()
				if err != nil {
					continue
				}

				data.Set(col.Name, result.Result[prop])
			default:
				result, err := ql.
					Data(props...).
					One()
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
* @param name string, page, rows int
* @return *Ql
**/
func (s *Ql) Detail(name string, page, rows int, tp TypeResult) *Ql {
	field := s.getField(name)
	if field == nil {
		return s
	}

	idx := slices.IndexFunc(s.Details, func(e *Field) bool { return e.AsField() == field.AsField() })
	if idx == -1 {
		s.Details = append(s.Details, field)
	}

	field.Page = page
	field.Rows = rows
	field.TpResult = tp

	return s
}

/**
* setDetail
* @param params []et.Json
* @return *Ql
**/
func (s *Ql) setDetail(params []et.Json) *Ql {
	for _, param := range params {
		name := param.Str("name")
		page := param.Int("page")
		rows := param.Int("rows")
		tp := StrToTypeResult(param.Str("type"))

		s.Detail(name, page, rows, tp)
	}

	return s
}

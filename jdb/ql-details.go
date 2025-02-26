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
			if col.GeneratedFunction != nil {
				col.GeneratedFunction(col, data)
			}
		case TpRelatedTo:
			if col.Detail == nil {
				continue
			}
			if col.Detail.Fk == nil {
				continue
			}
			pkn := col.Detail.Fk.Name
			key := (*data)[pkn]
			if key == nil {
				continue
			}

			fkn := col.Detail.Key
			with := col.Detail.With
			if field.TpResult == TpResult {
				result, err := with.
					Where(fkn).Eq(key).
					All()
				if err != nil {
					continue
				}

				data.Set(col.Name, result.Result)
			} else {
				all, err := with.
					Where(fkn).Eq(key).
					Counted()
				if err != nil {
					continue
				}

				result, err := with.
					Where(fkn).Eq(key).
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
			if col.Rollup.Fk == nil {
				continue
			}
			pkn := col.Rollup.Key
			key := (*data)[pkn]
			if key == nil {
				continue
			}

			n := len(col.Rollup.Props)
			if n <= 0 {
				continue
			}

			fkn := col.Rollup.Fk.Name
			source := col.Rollup.Source
			props := make([]string, 0)
			for _, prop := range col.Rollup.Props {
				props = append(props, prop.Name)
			}
			if n == 1 {
				prop := props[0]
				result, err := source.
					Where(fkn).Eq(key).
					Data(prop).
					One()
				if err != nil {
					continue
				}

				data.Set(col.Name, result.Result[prop])
			} else {
				result, err := source.
					Where(fkn).Eq(key).
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

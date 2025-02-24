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
	idx := slices.IndexFunc(s.Details, func(e *Field) bool { return e.Name == name })
	if idx == -1 {
		return s
	}

	field := s.Details[idx]
	if field == nil {
		return s
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

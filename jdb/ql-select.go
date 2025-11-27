package jdb

/**
* Select
* @param fields interface{}
* @return *Ql
**/
func (s *Ql) Select(fields ...string) *Ql {
	if len(s.Froms) == 0 {
		return s
	}

	for _, v := range fields {
		fld := &Field{
			Field: v,
		}
		fld = s.getField(fld)
		if !fld.IsDefined {
			continue
		}

		if TypeColumn[fld.Type] {
			s.Selects[fld.Field] = fld.As
			continue
		}

		if TypeAtrib[fld.Type] {
			s.Atribs[fld.Field] = fld.As
			continue
		}

		if fld.Type == TypeCalc {
			s.Calcs[v] = fld.Model.Calcs[fld.Name]
		} else if fld.Type == TypeDetail {
			s.Details[v] = fld.Model.Details[fld.Name]
		} else if fld.Type == TypeRollup {
			s.Rollups[v] = fld.Model.Rollups[fld.Name]
		} else if fld.Type == TypeRelation {
			s.Relations[v] = fld.Model.Relations[fld.Name]
		}
	}

	return s
}

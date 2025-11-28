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

	for _, name := range fields {
		fld := s.getField(name)
		if !fld.Existent {
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
			s.Calcs[fld.Name] = fld.Model.Calcs[fld.Name]
		} else if fld.Type == TypeDetail {
			s.Details[fld.Name] = fld.Model.Details[fld.Name]
		} else if fld.Type == TypeRollup {
			s.Rollups[fld.Name] = fld.Model.Rollups[fld.Name]
		} else if fld.Type == TypeRelation {
			s.Relations[fld.Name] = fld.Model.Relations[fld.Name]
		}
	}

	return s
}

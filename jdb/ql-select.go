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
			s.Selects[fld.As] = fld.Field
			continue
		}

		if TypeAtrib[fld.Type] {
			s.Atribs[fld.As] = fld.Field
			continue
		}

		if fld.Type == TypeCalc {
			s.Calcs[fld.As] = fld.Model.Calcs[fld.Name]
		} else if fld.Type == TypeDetail {
			s.Details[fld.As] = fld.Model.Details[fld.Name]
		} else if fld.Type == TypeRollup {
			s.Rollups[fld.As] = fld.Model.Rollups[fld.Name]
		} else if fld.Type == TypeRelation {
			s.Relations[fld.As] = fld.Model.Relations[fld.Name]
		}
	}

	return s
}

/**
* Data
* @param fields ...string
* @return *Ql
 */
func (s *Ql) Data(fields ...string) *Ql {
	s.IsDataSource = true
	return s.Select(fields...)
}

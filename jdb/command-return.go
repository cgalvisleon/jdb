package jdb

import "github.com/cgalvisleon/et/strs"

/**
* Return
* @param fields ...string
* @return *Command
**/
func (s *Command) Return(fields ...string) *Command {
	for _, name := range fields {
		sel := s.GetReturn(name)
		if sel != nil {
			s.Returns = append(s.Returns, sel)
		}
	}

	return s
}

/**
* listReturns
* @return []string
**/
func (s *Command) listReturns() []string {
	result := []string{}
	for _, sel := range s.Returns {
		result = append(result, strs.Format(`%s: %s`, sel.Field.AsField(), sel.Field.Name))
	}

	return result
}

package jdb

import (
	"slices"
)

/**
* before
* @return error
**/
func (s *Cmd) before() error {
	for _, item := range s.Data {
		for k, v := range item {
			_, ok := s.from.GetColumn(k)
			if !ok && !s.useAtribs {
				s.Atribs[k] = v
			}

			if !ok {
				continue
			}

			s.Columns[k] = v
			if slices.Contains(s.from.PrimaryKeys, k) {
				s.Keys[k] = v
			}
		}
	}

	return nil
}

/**
* after
* @return error
**/
func (s *Cmd) after() error {
	return nil
}

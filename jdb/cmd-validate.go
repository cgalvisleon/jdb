package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/et"
)

/**
* before
* @return error
**/
func (s *Cmd) before() error {
	for _, item := range s.Items {
		for k, v := range item {
			_, ok := s.from.GetColumn(k)
			if !ok && !s.useAtribs {
				s.Data[k] = et.Json{
					"type":  "atrib",
					"value": v,
				}
			}

			if !ok {
				continue
			}

			s.Data[k] = et.Json{
				"type":  "column",
				"value": v,
			}

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

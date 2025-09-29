package jdb

import "github.com/cgalvisleon/et/et"

/**
* validate
* @return error
**/
func (s *Model) validate() error {
	if len(s.Columns) != 0 {
		return nil
	}

	err := s.defineSourceField(SOURCE)
	if err != nil {
		return err
	}

	s.defineColumn(KEY, et.Json{
		"type": TypeKey,
	})
	s.definePrimaryKeys(KEY)
	s.defineRecordField(RECORDID)

	return nil
}

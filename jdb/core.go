package jdb

import "github.com/cgalvisleon/et/mistake"

var coreSchema *Schema

/**
* createCore
* @return error
**/
func (s *DB) createCore() error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}
	if err := s.defineModel(); err != nil {
		return err
	}
	if err := s.defineRecords(); err != nil {
		return err
	}
	if err := s.defineSeries(); err != nil {
		return err
	}
	if err := s.defineRecycling(); err != nil {
		return err
	}

	return nil
}

func (s *DB) defineSchema() error {
	if coreSchema != nil {
		return nil
	}

	var err error
	coreSchema, err = NewSchema(s, "core")
	if err != nil {
		return err
	}

	return nil
}

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
	if err := s.defineRecords(); err != nil {
		return err
	}
	if err := s.defineModel(); err != nil {
		return err
	}
	if err := s.defineSeries(); err != nil {
		return err
	}
	if err := s.defineRecycling(); err != nil {
		return err
	}

	s.isInit = true
	s.Save()
	coreSchema.Save()
	coreRecords.Save()
	coreModel.Save()
	coreSeries.Save()
	coreRecycling.Save()

	return nil
}

func (s *DB) defineSchema() error {
	if coreSchema != nil {
		return nil
	}

	coreSchema = NewSchema(s, "core")
	if coreSchema == nil {
		return mistake.New(MSG_SCHEMA_NOT_FOUND)
	}

	return nil
}

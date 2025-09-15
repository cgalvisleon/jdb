package jdb

import "fmt"

var coreSchema *Schema

/**
* createCore
* @return error
**/
func (s *DB) createCore() error {
	if s.driver == nil {
		return fmt.Errorf(MSG_DRIVER_NOT_DEFINED)
	}
	if err := s.defineModel(); err != nil {
		return err
	}
	if err := s.defineRecords(); err != nil {
		return err
	}
	if err := s.defineTables(); err != nil {
		return err
	}
	if err := s.defineRecycling(); err != nil {
		return err
	}
	if err := s.defineSeries(); err != nil {
		return err
	}

	return nil
}

func (s *DB) defineSchema() error {
	if coreSchema != nil {
		return nil
	}

	coreSchema = NewSchema(s, "core")
	if coreSchema == nil {
		return fmt.Errorf(MSG_SCHEMA_NOT_FOUND, "core")
	}

	return nil
}

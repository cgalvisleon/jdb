package postgres

import "github.com/cgalvisleon/et/console"

func (s *Postgres) defineCore() error {
	sql := parceSQL(`
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS pgcrypto;
	CREATE SCHEMA IF NOT EXISTS core;`)

	err := s.Exec(sql)
	if err != nil {
		return console.Panic(err)
	}

	return nil
}

func (s *Postgres) CreateCore() error {
	if err := s.defineCore(); err != nil {
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
	if err := s.defineDDL(); err != nil {
		return err
	}
	if err := s.defineModel(); err != nil {
		return err
	}
	if err := s.defineFunctions(s.nodeId); err != nil {
		return err
	}

	return nil
}

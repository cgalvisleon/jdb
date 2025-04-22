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

	return nil
}

package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/msg"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* LoadSchema
* @param name string
* @return error
**/
func (s *Postgres) LoadSchema(name string) error {
	if s.db == nil {
		return mistake.Newf(msg.NOT_DRIVER_DB)
	}

	exist, err := s.ExistSchema(name)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	sql := jdb.SQLDDL(`CREATE SCHEMA IF NOT EXISTS $1`, name)
	_, err = jdb.Query(s.db, sql)
	if err != nil {
		return err
	}

	console.Logf(s.name, `Schema %s created`, name)

	return nil
}

/**
* DropSchema
* @param name string
* @return error
**/
func (s *Postgres) DropSchema(name string) error {
	if s.db == nil {
		return mistake.Newf(msg.NOT_DRIVER_DB)
	}

	sql := jdb.SQLDDL(`DROP SCHEMA IF EXISTS $1 CASCADE`, name)
	_, err := jdb.Query(s.db, sql)
	if err != nil {
		return err
	}

	console.Logf(s.name, `Schema %s droped`, name)

	return nil
}

/**
* ExistSchema
* @param name string
* @return bool, error
**/
func (s *Postgres) ExistSchema(name string) (bool, error) {
	if s.db == nil {
		return false, mistake.Newf(msg.NOT_DRIVER_DB)
	}

	sql := jdb.SQLDDL(`SELECT 1 FROM pg_namespace WHERE nspname = '$1';`, name)
	items, err := jdb.Query(s.db, sql)
	if err != nil {
		return false, err
	}

	return items.Ok, nil
}

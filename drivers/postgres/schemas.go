package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/msg"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* loadSchema
* @param name string
* @return error
**/
func (s *Postgres) loadSchema(name string) error {
	if s.jdb == nil {
		return mistake.Newf(msg.NOT_DRIVER_DB)
	}

	exist, err := s.existSchema(name)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	sql := jdb.SQLDDL(`CREATE SCHEMA IF NOT EXISTS $1`, name)
	err = jdb.Ddl(s.jdb, sql)
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
	if s.jdb == nil {
		return mistake.Newf(msg.NOT_DRIVER_DB)
	}

	sql := jdb.SQLDDL(`DROP SCHEMA IF EXISTS $1 CASCADE`, name)
	err := jdb.Ddl(s.jdb, sql)
	if err != nil {
		return err
	}

	console.Logf(s.name, `Schema %s droped`, name)

	return nil
}

/**
* existSchema
* @param name string
* @return bool, error
**/
func (s *Postgres) existSchema(name string) (bool, error) {
	if s.jdb == nil {
		return false, mistake.Newf(msg.NOT_DRIVER_DB)
	}

	items, err := jdb.Query(s.jdb, `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.schemata
		WHERE UPPER(schema_name) = UPPER($1)
	);`, name)
	if err != nil {
		return false, err
	}

	return items.Ok, nil
}

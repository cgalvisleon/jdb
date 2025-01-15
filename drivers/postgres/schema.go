package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/msg"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) CreateSchema(name string) error {
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
	err = s.Exec(sql)
	if err != nil {
		return err
	}

	go s.upsertDDL(strs.Format(`create_schema_%s`, name), sql)

	console.Logf(jdb.Postgres, `Schema %s created`, name)

	return nil
}

func (s *Postgres) DropSchema(name string) error {
	if s.db == nil {
		return mistake.Newf(msg.NOT_DRIVER_DB)
	}

	sql := jdb.SQLDDL(`DROP SCHEMA IF EXISTS $1 CASCADE`, name)
	err := s.Exec(sql)
	if err != nil {
		return err
	}

	go s.upsertDDL(strs.Format(`drop_schema_%s`, name), sql)

	console.Logf(jdb.Postgres, `Schema %s droped`, name)

	return nil
}

func (s *Postgres) ExistSchema(name string) (bool, error) {
	if s.db == nil {
		return false, mistake.Newf(msg.NOT_DRIVER_DB)
	}

	sql := jdb.SQLDDL(`SELECT 1 FROM pg_namespace WHERE nspname = '$1';`, name)
	items, err := s.SQL(sql)
	if err != nil {
		return false, err
	}

	return items.Ok, nil
}

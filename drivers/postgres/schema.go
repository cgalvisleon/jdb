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

	sql := jdb.SQLDDL(`CREATE SCHEMA IF NOT EXISTS $1`, name)
	err := s.Exec(sql)
	if err != nil {
		return mistake.Newf(jdb.MSG_QUERY_FAILED, err.Error())
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
		return mistake.Newf(jdb.MSG_QUERY_FAILED, err.Error())
	}

	go s.upsertDDL(strs.Format(`drop_schema_%s`, name), sql)

	console.Logf(jdb.Postgres, `Schema %s droped`, name)

	return nil
}

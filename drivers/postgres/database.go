package postgres

import (
	"github.com/cgalvisl/jdb/jdb"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/msg"
	"github.com/cgalvisleon/et/strs"
)

func (s *Postgres) existDatabase(name string) (bool, error) {
	sql := strs.Format(`SELECT 1 FROM pg_database WHERE datname = %s;`, name)
	items, err := s.SQL(sql)
	if err != nil {
		return false, err
	}

	if !items.Ok {
		return false, logs.Alertm(jdb.MSG_DATABASE_NOT_FOUND)
	}

	return true, nil
}

func (s *Postgres) CreateDatabase(name string) error {
	if s.db == nil {
		return logs.Alertf(msg.NOT_DRIVER_DB)
	}

	exist, err := s.existDatabase(name)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	sql := strs.Format(`CREATE DATABASE `, name)
	err = s.Exec(sql)
	if err != nil {
		return logs.Alertf(jdb.MSG_QUERY_FAILED, err.Error())
	}

	return nil
}

func (s *Postgres) DropDatabase(name string) error {
	if s.db == nil {
		return logs.Alertf(msg.NOT_DRIVER_DB)
	}

	exist, err := s.existDatabase(name)
	if err != nil {
		return err
	}

	if !exist {
		return nil
	}

	sql := strs.Format(`DROP DATABASE `, name)
	err = s.Exec(sql)
	if err != nil {
		return logs.Alertf(jdb.MSG_QUERY_FAILED, err.Error())
	}

	return nil
}

func (s *Postgres) RenameDatabase(name, newname string) error {
	return nil
}

func (s *Postgres) SetParams(data et.Json) error {
	return nil
}

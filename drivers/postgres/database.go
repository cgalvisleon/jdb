package postgres

import (
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/msg"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) existDatabase(name string) (bool, error) {
	sql := `
	SELECT EXISTS(
	SELECT 1
	FROM pg_database
	WHERE UPPER(datname) = UPPER($1));`
	items, err := s.SQL(sql, name)
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

	sql := `CREATE DATABASE $1`
	err = s.Exec(sql, name)
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

	sql := strs.Format(`DROP DATABASE %s`, name)
	err = s.Exec(sql)
	if err != nil {
		return logs.Alertf(jdb.MSG_QUERY_FAILED, err.Error())
	}

	return nil
}

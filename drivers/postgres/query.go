package postgres

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	jdb "github.com/cgalvisleon/jdb/pkg"
)

func (s *Postgres) query(db *sql.DB, sql string, params ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(sql, params...)
	if err != nil {
		return nil, logs.Alertf(jdb.MSG_QUERY_FAILED, err.Error())
	}

	return rows, nil
}

func (s *Postgres) Exec(sql string, params ...interface{}) error {
	_, err := s.query(s.db, sql, params...)
	if err != nil {
		return err
	}

	return nil

}

func (s *Postgres) SQL(sql string, params ...interface{}) (et.Items, error) {
	rows, err := s.query(s.db, sql, params...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	result := jdb.RowsToItems(rows)

	return result, nil
}

func (s *Postgres) One(sql string, params ...interface{}) (et.Item, error) {
	rows, err := s.query(s.db, sql, params...)
	if err != nil {
		return et.Item{}, err
	}
	defer rows.Close()

	result := jdb.RowsToItem(rows)

	return result, nil
}

func (s *Postgres) Query(linq *jdb.Linq) (et.Items, error) {
	return et.Items{}, nil
}

func (s *Postgres) Count(linq *jdb.Linq) (int, error) {
	return 0, nil
}

func (s *Postgres) Last(linq *jdb.Linq) (et.Items, error) {
	return et.Items{}, nil
}

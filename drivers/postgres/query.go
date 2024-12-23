package postgres

import (
	"database/sql"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) query(db *sql.DB, sql string, params ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(sql, params...)
	if err != nil {
		sql = jdb.SQLParse(sql, params...)
		return nil, console.QueryError(err, sql)
	}

	return rows, nil
}

func (s *Postgres) exec(db *sql.DB, sql string, params ...interface{}) (sql.Result, error) {
	result, err := db.Exec(sql, params...)
	if err != nil {
		sql = jdb.SQLParse(sql, params...)
		return nil, console.QueryError(err, sql)
	}

	return result, nil
}

func (s *Postgres) Exec(sql string, params ...interface{}) error {
	_, err := s.exec(s.db, sql, params...)
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

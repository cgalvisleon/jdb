package postgres

import (
	"database/sql"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) query(db *sql.DB, sql string, arg ...any) (*sql.Rows, error) {
	result, err := db.Query(sql, arg...)
	if err != nil {
		sql = jdb.SQLParse(sql, arg...)
		return nil, console.QueryError(err, sql)
	}

	return result, nil
}

func (s *Postgres) exec(db *sql.DB, sql string, arg ...any) (sql.Result, error) {
	result, err := db.Exec(sql, arg...)
	if err != nil {
		sql = jdb.SQLParse(sql, arg...)
		return nil, console.QueryError(err, sql)
	}

	return result, nil
}

/**
* Exec
* @param sql string, arg ...any
* @return error
**/
func (s *Postgres) Exec(sql string, arg ...any) error {
	_, err := s.exec(s.db, sql, arg...)
	if err != nil {
		return err
	}

	return nil

}

/**
* QueryRow
* @param query string, dest ...any
* @return error
**/
func (s *Postgres) QueryRow(query string, dest ...any) (bool, error) {
	err := s.db.QueryRow(query).Scan(dest...)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, err
	}

	return false, nil
}

/**
* Query
* @param sql string, arg ...any
* @return et.Items, error
**/
func (s *Postgres) Query(sql string, arg ...any) (et.Items, error) {
	rows, err := s.query(s.db, sql, arg...)
	if err != nil {
		return et.Items{}, console.QueryError(err, sql)
	}
	defer rows.Close()

	result := jdb.RowsToItems(rows)

	return result, nil
}

/**
* One
* @param sql string, arg ...any
* @return et.Item, error
**/
func (s *Postgres) One(sql string, arg ...any) (et.Item, error) {
	rows, err := s.query(s.db, sql, arg...)
	if err != nil {
		return et.Item{}, console.QueryError(err, sql)
	}
	defer rows.Close()

	result := jdb.RowsToItem(rows)

	return result, nil
}

/**
* Data
* @param source, sql string, arg ...any
* @return et.Items, error
**/
func (s *Postgres) Data(sourceFiled, sql string, arg ...any) (et.Items, error) {
	rows, err := s.query(s.db, sql, arg...)
	if err != nil {
		return et.Items{}, console.QueryError(err, sql)
	}
	defer rows.Close()

	result := jdb.RowsToSource(sourceFiled, rows)

	return result, nil
}

/**
* Select
* @param ql *jdb.Ql
* @return et.Items, error
**/
func (s *Postgres) Select(ql *jdb.Ql) (et.Items, error) {
	ql.Sql = ""
	ql.Sql = strs.Append(ql.Sql, s.sqlSelect(ql), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlFrom(ql.Froms), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlJoin(ql.Joins), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlWhere(ql.QlWhere), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlGroupBy(ql), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlHaving(ql), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlOrderBy(ql), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlLimit(ql), "\n")
	ql.Sql = strs.Format(`%s;`, ql.Sql)

	if ql.IsDebug {
		console.Debug(ql.Sql)
	}

	if ql.TypeSelect == jdb.Data {
		result, err := s.Data("result", ql.Sql)
		if err != nil {
			return et.Items{}, err
		}

		return result, nil
	}

	result, err := s.Query(ql.Sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* ExecDDL
* @param id, sql string, arg ...any
* @return error
**/
func (s *Postgres) ExecDDL(id, sql string, arg ...any) error {
	err := s.Exec(sql, arg...)
	if err != nil {
		return err
	}

	go s.upsertDDL(id, sql)

	return nil
}

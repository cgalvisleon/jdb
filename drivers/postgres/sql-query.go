package postgres

import (
	"database/sql"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) query(db *sql.DB, sql string, params ...any) (*sql.Rows, error) {
	result, err := db.Query(sql, params...)
	if err != nil {
		sql = jdb.SQLParse(sql, params...)
		return nil, console.QueryError(err, sql)
	}

	return result, nil
}

func (s *Postgres) exec(db *sql.DB, sql string, params ...any) (sql.Result, error) {
	result, err := db.Exec(sql, params...)
	if err != nil {
		sql = jdb.SQLParse(sql, params...)
		return nil, console.QueryError(err, sql)
	}

	return result, nil
}

/**
* Exec
* @param sql string, params ...any
* @return error
**/
func (s *Postgres) Exec(sql string, params ...any) error {
	_, err := s.exec(s.db, sql, params...)
	if err != nil {
		return err
	}

	return nil

}

/**
* Query
* @param sql string, params ...any
* @return et.Items, error
**/
func (s *Postgres) Query(sql string, params ...any) (et.Items, error) {
	rows, err := s.query(s.db, sql, params...)
	if err != nil {
		return et.Items{}, console.QueryError(err, sql)
	}
	defer rows.Close()

	result := jdb.RowsToItems(rows)

	return result, nil
}

/**
* Data
* @param source, sql string, params ...any
* @return et.Items, error
**/
func (s *Postgres) Data(sourceFiled, sql string, params ...any) (et.Items, error) {
	rows, err := s.query(s.db, sql, params...)
	if err != nil {
		return et.Items{}, console.QueryError(err, sql)
	}
	defer rows.Close()

	result := jdb.SourceToItems(sourceFiled, rows)

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

	for _, item := range result.Result {
		for _, detail := range ql.Details {
			item[detail.Name] = detail.Define()
			console.Debug(detail)
		}
	}

	return result, nil
}

/**
* ExecDDL
* @param id, sql string, params ...any
* @return error
**/
func (s *Postgres) ExecDDL(id, sql string, params ...any) error {
	err := s.Exec(sql, params...)
	if err != nil {
		return err
	}

	go s.upsertDDL(id, sql)

	return nil
}

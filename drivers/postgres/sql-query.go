package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* queryTx
* @param tx *sql.Tx, sql string, arg ...any
* @return *sql.Rows, error
**/
func (s *Postgres) queryTx(tx *jdb.Tx, sql string, arg ...any) (et.Items, error) {
	if tx != nil {
		err := tx.Begin(s.db)
		if err != nil {
			return et.Items{}, err
		}

		rows, err := tx.Tx.Query(sql, arg...)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				err = fmt.Errorf("error on rollback: %w: %s", errRollback, err)
			}

			return et.Items{}, console.QueryError(err, sql)
		}
		defer rows.Close()

		result := jdb.RowsToItems(rows)

		return result, nil
	}

	rows, err := s.db.Query(sql, arg...)
	if err != nil {
		sql = jdb.SQLParse(sql, arg...)
		return et.Items{}, console.QueryError(err, sql)
	}
	defer rows.Close()

	result := jdb.RowsToItems(rows)

	return result, nil
}

/**
* query
* @param sql string, arg ...any
* @return et.Items, error
**/
func (s *Postgres) query(sql string, arg ...any) (et.Items, error) {
	return s.queryTx(nil, sql, arg...)
}

/**
* data
* @param tx *sql.Tx, sourceFiled, sql string, arg ...any
* @return et.Items, error
**/
func (s *Postgres) dataTx(tx *jdb.Tx, sourceFiled, sql string, arg ...any) (et.Items, error) {
	if tx != nil {
		err := tx.Begin(s.db)
		if err != nil {
			return et.Items{}, err
		}

		rows, err := tx.Tx.Query(sql, arg...)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				err = fmt.Errorf("error on rollback: %w: %s", errRollback, err)
			}

			return et.Items{}, console.QueryError(err, sql)
		}
		defer rows.Close()

		result := jdb.RowsToSource(sourceFiled, rows)

		return result, nil
	}

	rows, err := s.db.Query(sql, arg...)
	if err != nil {
		sql = jdb.SQLParse(sql, arg...)
		return et.Items{}, console.QueryError(err, sql)
	}
	defer rows.Close()

	result := jdb.RowsToSource(sourceFiled, rows)

	return result, nil
}

/**
* QueryTx
* @param tx *jdb.Tx, sql string, arg ...any
* @return et.Items, error
**/
func (s *Postgres) QueryTx(tx *jdb.Tx, sql string, arg ...any) (et.Items, error) {
	return s.queryTx(tx, sql, arg...)
}

/**
* Query
* @param sql string, arg ...any
* @return et.Items, error
**/
func (s *Postgres) Query(sql string, arg ...any) (et.Items, error) {
	return s.query(sql, arg...)
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
		result, err := s.dataTx(ql.Tx(), "result", ql.Sql)
		if err != nil {
			return et.Items{}, err
		}

		return result, nil
	}

	result, err := s.queryTx(ql.Tx(), ql.Sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

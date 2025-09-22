package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
)

/**
* querytx
* @param db *sql.DB, tx *Tx, sql string, arg ...any
* @return *sql.Rows, error
**/
func querytx(db *Database, tx *Tx, sql string, arg ...any) (et.Items, error) {
	sql = SQLParse(sql, arg...)
	data := et.Json{
		"db_name": db.Name,
		"sql":     sql,
	}

	if tx != nil {
		err := tx.Begin(db.Db)
		if err != nil {
			return et.Items{}, err
		}

		rows, err := tx.Tx.Query(sql, arg...)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				data["error"] = err.Error()
				event.Publish(EVENT_SQL_ERROR, data)
				err = fmt.Errorf("error on rollback: %w: %s", errRollback, err)
			}

			return et.Items{}, err
		}
		defer rows.Close()

		if tx.Committed {
			tp := tipoSQL(sql)
			event.Publish(fmt.Sprintf("sql:%s", tp), data)
		}

		return RowsToItems(rows), nil
	}

	rows, err := db.Db.Query(sql, arg...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	tp := tipoSQL(sql)
	event.Publish(fmt.Sprintf("sql:%s", tp), data)
	return RowsToItems(rows), nil
}

/**
* QueryTx
* @param tx *Tx, sql string, arg ...any
* @return et.Items, error
**/
func (s *Database) QueryTx(tx *Tx, sql string, arg ...any) (et.Items, error) {
	return querytx(s, tx, sql, arg...)
}

/**
* Query
* @param sql string, arg ...any
* @return et.Items, error
**/
func (s *Database) Query(sql string, arg ...any) (et.Items, error) {
	return querytx(s, nil, sql, arg...)
}

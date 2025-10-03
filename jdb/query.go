package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
)

/**
* querytx
* @param db *sql.DB, tx *Tx, hidden []string, sql string, arg ...any
* @return *sql.Rows, error
**/
func querytx(db *DB, tx *Tx, hidden []string, sql string, arg ...any) (et.Items, error) {
	data := et.Json{
		"db_name": db.Name,
		"sql":     sql,
		"args":    arg,
		"hidden":  hidden,
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

		return RowsToItems(rows, hidden), nil
	}

	rows, err := db.Db.Query(sql, arg...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	tp := tipoSQL(sql)
	event.Publish(fmt.Sprintf("sql:%s", tp), data)
	return RowsToItems(rows, hidden), nil
}

/**
* QueryTx
* @param tx *Tx, hidden []string, sql string, arg ...any
* @return et.Items, error
**/
func (s *DB) QueryTx(tx *Tx, hidden []string, sql string, arg ...any) (et.Items, error) {
	return querytx(s, tx, hidden, sql, arg...)
}

/**
* Query
* @param hidden []string, sql string, arg ...any
* @return et.Items, error
**/
func (s *DB) Query(hidden []string, sql string, arg ...any) (et.Items, error) {
	return querytx(s, nil, hidden, sql, arg...)
}

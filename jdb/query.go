package jdb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
)

func TipoSQL(query string) string {
	q := strings.TrimSpace(strings.ToUpper(query))

	parts := strings.Fields(q)
	if len(parts) == 0 {
		return "DESCONOCIDO"
	}

	cmd := parts[0]

	switch cmd {
	case "SELECT":
		return "query"
	case "INSERT", "UPDATE", "DELETE", "MERGE":
		return "command"
	case "CREATE", "ALTER", "DROP", "TRUNCATE":
		return "definition"
	case "GRANT", "REVOKE":
		return "definition"
	case "COMMIT", "ROLLBACK", "SAVEPOINT", "SET":
		return "definition"
	default:
		return "desconocido"
	}
}

/**
* queryTx
* @param db *sql.DB, tx *Tx, sourceFiled, sql string, arg ...any
* @return *sql.Rows, error
**/
func queryTx(db *sql.DB, tx *Tx, sourceFiled, sql string, arg ...any) (et.Items, error) {
	if tx != nil {
		err := tx.Begin(db)
		if err != nil {
			return et.Items{}, err
		}

		rows, err := tx.Tx.Query(sql, arg...)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				err = fmt.Errorf("error on rollback: %w: %s", errRollback, err)
			}

			return et.Items{}, fmt.Errorf("QueryTx error: %s", err.Error())
		}
		defer rows.Close()

		if sourceFiled != "" {
			return RowsToSourceItems(rows, sourceFiled), nil
		}

		return RowsToItems(rows), nil
	}

	rows, err := db.Query(sql, arg...)
	if err != nil {
		return et.Items{}, fmt.Errorf("Query error: %s", err.Error())
	}
	defer rows.Close()

	if sourceFiled != "" {
		return RowsToSourceItems(rows, sourceFiled), nil
	}

	return RowsToItems(rows), nil
}

/**
* query
* @param db *DB, tx *Tx, sourceFiled, sql string, arg ...any
* @return et.Items, error
**/
func query(db *DB, tx *Tx, sourceFiled, sql string, arg ...any) (et.Items, error) {
	result, err := queryTx(db.db, tx, sourceFiled, sql, arg...)
	if err != nil {
		event.Publish(EVENT_SQL_ERROR, et.Json{
			"db_name": db.Name,
			"sql":     sql,
			"arg":     arg,
			"error":   err,
		})
		return et.Items{}, err
	}

	tp := TipoSQL(sql)
	event.Publish(fmt.Sprintf("sql:%s", tp), et.Json{
		"db_name": db.Name,
		"sql":     sql,
		"arg":     arg,
	})

	return result, nil
}

/**
* QueryTx
* @param db *DB, tx *Tx, sourceFiled, sql string, arg ...any
* @return et.Items, error
**/
func QueryTx(db *DB, tx *Tx, sql string, arg ...any) (et.Items, error) {
	return query(db, tx, "", sql, arg...)
}

/**
* Query
* @param db *DB, sql string, arg ...any
* @return et.Items, error
**/
func Query(db *DB, sql string, arg ...any) (et.Items, error) {
	return QueryTx(db, nil, sql, arg...)
}

/**
* ResultTx
* @param db *DB, tx *Tx, sql string, arg ...any
* @return et.Items, error
**/
func ResultTx(db *DB, tx *Tx, sql string, arg ...any) (et.Items, error) {
	return query(db, tx, "result", sql, arg...)
}

/**
* Result
* @param db *DB, tx *Tx, sql string, arg ...any
* @return et.Items, error
**/
func Result(db *DB, sql string, arg ...any) (et.Items, error) {
	return ResultTx(db, nil, sql, arg...)
}

/**
* DefinitionTx
* @param db *DB, tx *Tx, sql string, arg ...any
* @return error
**/
func DefinitionTx(db *DB, tx *Tx, sql string, arg ...any) error {
	_, err := query(db, tx, "", sql, arg...)
	if err != nil {
		return err
	}

	return nil
}

/**
* Definition
* @param db *DB, sql string, arg ...any
* @return error
**/
func Definition(db *DB, sql string, arg ...any) error {
	return DefinitionTx(db, nil, sql, arg...)
}

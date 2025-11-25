package jdb

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
)

const (
	COMMAND     = "command"
	QUERY       = "query"
	DEFINITION  = "definition"
	CONTROL     = "control"
	TRANSACTION = "transaction"
	OTHER       = "other"
)

/**
* tipoSQL
* @param query string
* @return string
**/
func tipoSQL(query string) string {
	q := strings.TrimSpace(strings.ToUpper(query))

	parts := strings.Fields(q)
	if len(parts) == 0 {
		return OTHER
	}

	cmd := parts[0]

	switch cmd {
	case "SELECT":
		return QUERY
	case "INSERT", "UPDATE", "DELETE", "MERGE":
		return COMMAND
	case "CREATE", "ALTER", "DROP", "TRUNCATE":
		return DEFINITION
	case "GRANT", "REVOKE":
		return CONTROL
	case "COMMIT", "ROLLBACK", "SAVEPOINT", "SET":
		return TRANSACTION
	default:
		return OTHER
	}
}

/**
* RowsToItems
* @param rows *sql.Rows
* @return et.Items
**/
func RowsToItems(rows *sql.Rows) et.Items {
	var result = et.Items{Result: []et.Json{}}

	append := func(item et.Json) {
		result.Add(item)
	}

	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		if len(item) == 1 {
			for _, v := range item {
				switch val := v.(type) {
				case et.Json:
					append(val)
				case map[string]interface{}:
					append(et.Json(val))
				default:
					append(item)
				}
			}
		} else {
			append(item)
		}
	}

	return result
}

/**
* SQLUnQuote
* @param sql string
* @param args ...any
* @return string
**/
func SQLUnQuote(sql string, args ...any) string {
	for i := range args {
		old := fmt.Sprintf(`$%d`, i+1)
		new := fmt.Sprintf(`{$%d}`, i+1)
		sql = strings.ReplaceAll(sql, old, new)
	}

	for i, arg := range args {
		old := fmt.Sprintf(`{$%d}`, i+1)
		new := fmt.Sprintf(`%v`, arg)
		sql = strings.ReplaceAll(sql, old, new)
	}

	return sql
}

/**
* SQLParse
* @param sql string
* @param args ...any
* @return string
**/
func SQLParse(sql string, args ...any) string {
	for i := range args {
		old := fmt.Sprintf(`$%d`, i+1)
		new := fmt.Sprintf(`{$%d}`, i+1)
		sql = strings.ReplaceAll(sql, old, new)
	}

	for i, arg := range args {
		old := fmt.Sprintf(`{$%d}`, i+1)
		new := fmt.Sprintf(`%v`, Quote(arg))
		sql = strings.ReplaceAll(sql, old, new)
	}

	return sql
}

/**
* Quote
* @param val any
* @return any
**/
func Quote(val any) any {
	format := `'%v'`
	switch v := val.(type) {
	case string:
		return fmt.Sprintf(format, v)
	case int:
		return v
	case float64:
		return v
	case float32:
		return v
	case int16:
		return v
	case int32:
		return v
	case int64:
		return v
	case bool:
		return v
	case time.Time:
		return fmt.Sprintf(format, v.Format("2006-01-02 15:04:05"))
	case et.Json:
		return fmt.Sprintf(format, v.ToString())
	case map[string]interface{}:
		return fmt.Sprintf(format, et.Json(v).ToString())
	case []string, []et.Json, []interface{}, []map[string]interface{}:
		bt, err := json.Marshal(v)
		if err != nil {
			logs.Errorf("Quote, type:%v, value:%v, error marshalling array: %v", reflect.TypeOf(v), v, err)
			return strs.Format(format, `[]`)
		}
		return fmt.Sprintf(format, string(bt))
	case []uint8:
		b := []byte(val.([]uint8))
		return fmt.Sprintf("'\\x%s'", hex.EncodeToString(b))
	case nil:
		return fmt.Sprintf(`%s`, "NULL")
	default:
		logs.Errorf("Quote, type:%v, value:%v", reflect.TypeOf(v), v)
		return val
	}
}

/**
* GetFieldName
* @param s string
* @return string
**/
func GetFieldName(s string) string {
	result := s
	if strings.Contains(s, ".") {
		parts := strings.Split(s, ".")
		result = parts[len(parts)-1]
	}
	return result
}

/**
* GetAtribName
* @param s string
* @return string
**/
func GetAtribName(s string) string {
	result := GetFieldName(s)
	if strings.Contains(result, ":") {
		parts := strings.Split(result, ":")
		result = parts[len(parts)-1]
	}
	return result
}

/**
* querytx
* @param db *sql.DB, tx *Tx, query string, arg ...any
* @return *sql.Rows, error
**/
func querytx(db *DB, tx *Tx, query string, arg ...any) (et.Items, error) {
	data := et.Json{
		"db_name": db.Name,
		"query":   query,
		"args":    arg,
	}

	var err error
	var rows *sql.Rows
	if tx != nil {
		err = tx.Begin(db.Db)
		if err != nil {
			return et.Items{}, err
		}

		rows, err = tx.Tx.Query(query, arg...)
		if err != nil {
			err = fmt.Errorf(`%s: %w`, query, err)
			data["error"] = err.Error()
			event.Publish(EVENT_SQL_ERROR, data)
			errRollback := tx.Rollback()
			if errRollback != nil {
				err = fmt.Errorf("error on rollback: %w: %s", errRollback, err)
			}

			return et.Items{}, err
		}
	} else {
		rows, err = db.Db.Query(query, arg...)
		if err != nil {
			err = fmt.Errorf(`%s: %w`, query, err)
			data["error"] = err.Error()
			event.Publish(EVENT_SQL_ERROR, data)
			return et.Items{}, err
		}
	}

	tp := tipoSQL(query)
	event.Publish(fmt.Sprintf("sql:%s", tp), data)
	defer rows.Close()
	result := RowsToItems(rows)
	return result, nil
}

/**
* QueryTx
* @param tx *Tx, sql string, arg ...any
* @return et.Items, error
**/
func (s *DB) QueryTx(tx *Tx, sql string, arg ...any) (et.Items, error) {
	return querytx(s, tx, sql, arg...)
}

/**
* Query
* @param sql string, arg ...any
* @return et.Items, error
**/
func (s *DB) Query(sql string, arg ...any) (et.Items, error) {
	return querytx(s, nil, sql, arg...)
}

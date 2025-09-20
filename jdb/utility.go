package jdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
)

const (
	COMMAND    = "command"
	QUERY      = "query"
	DEFINITION = "definition"
	STRANGE    = "strange"
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
		return STRANGE
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
		return DEFINITION
	case "COMMIT", "ROLLBACK", "SAVEPOINT", "SET":
		return DEFINITION
	default:
		return STRANGE
	}
}

/**
* RowsToItems
* @param rows *sql.Rows
* @return et.Items
**/
func RowsToItems(rows *sql.Rows) et.Items {
	var result = et.Items{Result: []et.Json{}}
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Count++
		if len(item) == 1 {
			for _, v := range item {
				switch val := v.(type) {
				case et.Json:
					result.Result = append(result.Result, val)
				case map[string]interface{}:
					result.Result = append(result.Result, et.Json(val))
				default:
					result.Result = append(result.Result, item)
				}
			}
		} else {
			result.Result = append(result.Result, item)
		}
	}

	return result
}

/**
* rowsToSourceItems
* @param rows *sql.Rows, source string
* @return et.Items
**/
func rowsToSourceItems(rows *sql.Rows, source string) et.Items {
	var result = et.Items{Result: []et.Json{}}
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Count++
		if item[source] == nil {
			result.Result = append(result.Result, item)
		} else {
			result.Result = append(result.Result, item.Json(source))
		}
	}

	return result
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
* @param val interface{}
* @return any
**/
func Quote(val interface{}) any {
	fmt := `'%s'`
	switch v := val.(type) {
	case string:
		return strconv.Quote(v)
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
		return strs.Format(fmt, v.Format("2006-01-02 15:04:05"))
	case et.Json:
		return strs.Format(fmt, v.ToString())
	case map[string]interface{}:
		return strs.Format(fmt, et.Json(v).ToString())
	case []string, []et.Json, []interface{}, []map[string]interface{}:
		bt, err := json.Marshal(v)
		if err != nil {
			logs.Errorf("Quote", "type:%v, value:%v, error marshalling array: %v", reflect.TypeOf(v), v, err)
			return strs.Format(fmt, `[]`)
		}
		return strs.Format(fmt, string(bt))
	case []uint8:
		return strs.Format(fmt, string(v))
	case nil:
		return strs.Format(`%s`, "NULL")
	default:
		logs.Errorf("Quote", "type:%v, value:%v", reflect.TypeOf(v), v)
		return val
	}
}

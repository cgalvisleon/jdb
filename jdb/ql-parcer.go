package jdb

import (
	"database/sql"
	"reflect"
	"strings"
	"time"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

/**
* SQLQuote return a sql cuote string
* @param sql string
* @return string
**/
func SQLQuote(sql string) string {
	sql = strings.TrimSpace(sql)

	result := strs.Replace(sql, `'`, `"`)
	result = strs.Trim(result)

	return result
}

/**
* SQLDDL return a sql string with the args
* @param sql string
* @param args ...any
* @return string
**/
func SQLDDL(sql string, args ...any) string {
	sql = strings.TrimSpace(sql)

	for i, arg := range args {
		old := strs.Format(`$%d`, i+1)
		new := strs.Format(`%v`, arg)
		sql = strings.ReplaceAll(sql, old, new)
	}

	return sql
}

/**
* SQLParse return a sql string with the args
* @param sql string
* @param args ...any
* @return string
**/
func SQLParse(sql string, args ...any) string {
	for i := range args {
		old := strs.Format(`$%d`, i+1)
		new := strs.Format(`{$%d}`, i+1)
		sql = strings.ReplaceAll(sql, old, new)
	}

	for i, arg := range args {
		old := strs.Format(`{$%d}`, i+1)
		new := strs.Format(`%v`, Quote(arg))
		sql = strings.ReplaceAll(sql, old, new)
	}

	return sql
}

/**
* RowsToItems return a items from a sql rows
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
		result.Result = append(result.Result, item)
	}

	return result
}

/**
* RowsToItem return a item from a sql rows
* @param rows *sql.Rows
* @return et.Item
**/
func RowsToItem(rows *sql.Rows) et.Item {
	var result = et.Item{Result: et.Json{}}
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Result = item
		break
	}

	return result
}

/**
* SourceToItems return a items from a sql rows and source field
* @param source string
* @param rows *sql.Rows
* @return et.Items
**/
func RowsToSource(source string, rows *sql.Rows) et.Items {
	var result = et.Items{Result: []et.Json{}}
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Count++
		result.Result = append(result.Result, item.Json(source))
	}

	return result
}

/**
* SourceToItem return a item from a sql rows and source field
* @param source string
* @param rows *sql.Rows
* @return et.Items
**/
func SourceToItem(source string, rows *sql.Rows) et.Item {
	var result = et.Item{Result: et.Json{}}
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Result = item.Json(source)
		break
	}

	return result
}

/**
* JsonQuote return a json quote string
* @param val interface{}
* @return interface{}
**/
func JsonQuote(val interface{}) interface{} {
	fmt := `'%v'`
	switch v := val.(type) {
	case string:
		v = strs.Format(`"%s"`, v)
		return strs.Format(fmt, v)
	case int:
		return strs.Format(fmt, v)
	case float64:
		return strs.Format(fmt, v)
	case float32:
		return strs.Format(fmt, v)
	case int16:
		return strs.Format(fmt, v)
	case int32:
		return strs.Format(fmt, v)
	case int64:
		return strs.Format(fmt, v)
	case bool:
		return strs.Format(fmt, v)
	case time.Time:
		return strs.Format(fmt, v.Format("2006-01-02 15:04:05"))
	case et.Json:
		return strs.Format(fmt, v.ToString())
	case map[string]interface{}:
		return strs.Format(fmt, et.Json(v).ToString())
	case []string:
		var r string
		for _, s := range v {
			r = strs.Append(r, strs.Format(`"%s"`, s), ", ")
		}
		r = strs.Format(`[%s]`, r)
		return strs.Format(fmt, r)
	case []interface{}:
		var r string
		for _, _v := range v {
			q := JsonQuote(_v)
			r = strs.Append(r, strs.Format(`%v`, q), ", ")
		}
		r = strs.Format(`[%s]`, r)
		return strs.Format(fmt, r)
	case []et.Json:
		var r string
		for _, _v := range v {
			q := JsonQuote(_v)
			r = strs.Append(r, strs.Format(`%v`, q), ", ")
		}
		r = strs.Format(`[%s]`, r)
		return strs.Format(fmt, r)
	case []map[string]interface{}:
		var r string
		for _, _v := range v {
			q := JsonQuote(_v)
			r = strs.Append(r, strs.Format(`%v`, q), ", ")
		}
		r = strs.Format(`[%s]`, r)
		return strs.Format(fmt, r)
	case []uint8:
		return strs.Format(fmt, string(v))
	case nil:
		return strs.Format(`%s`, "NULL")
	default:
		console.Errorf("JsonQuote type:%v value:%v", reflect.TypeOf(v), v)
		return val
	}
}

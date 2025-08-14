package jdb

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
)

var quotedChar = `'`

/**
* SetQuotedChar
* @param char string
**/
func SetQuotedChar(char string) {
	quotedChar = strs.Format(`%s`, char)
}

/**
* unquote
* @param str string
* @return string
**/
func unquote(str string) string {
	str = strings.ReplaceAll(str, `'`, `"`)
	result, err := strconv.Unquote(str)
	if err != nil {
		result = str
	}

	return result
}

/**
* quote
* @param str string
* @return string
**/
func quote(str string) string {
	result := strconv.Quote(str)
	if quotedChar == `"` {
		return result
	}

	return strings.ReplaceAll(result, `"`, `'`)
}

/**
* Unquote
* @param val interface{}
* @return any
**/
func Unquote(val interface{}) any {
	switch v := val.(type) {
	case string:
		return unquote(v)
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
	case et.Json:
		return strs.Format(`%s`, v.ToString())
	case map[string]interface{}:
		return strs.Format(`%s`, et.Json(v).ToString())
	case time.Time:
		return strs.Format(`%s`, v.Format("2006-01-02 15:04:05"))
	case []string:
		var r string
		for i, _v := range v {
			if i == 0 {
				r = strs.Format(`%s`, unquote(_v))
			} else {
				r = strs.Format(`%s, %s`, r, unquote(_v))
			}
		}
		return strs.Format(`[%s]`, unquote(r))
	case []interface{}:
		var r string
		for i, _v := range v {
			q := Unquote(_v)
			if i == 0 {
				r = strs.Format(`%v`, q)
			} else {
				r = strs.Format(`%s, %v`, r, q)
			}
		}
		return strs.Format(`[%s]`, r)
	case []uint8:
		return strs.Format(`%s`, string(v))
	case nil:
		return strs.Format(`%s`, "NULL")
	default:
		logs.Errorf("Not unquoted type:%v value:%v", reflect.TypeOf(v), v)
		return val
	}
}

/**
* Quote
* @param val interface{}
* @return any
**/
func Quote(val interface{}) any {
	fmt := `'%s'`
	if quotedChar == `"` {
		fmt = `"%s"`
	}
	switch v := val.(type) {
	case string:
		return quote(v)
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
	case []et.Json, []string, []interface{}, []map[string]interface{}:
		bt, err := json.Marshal(v)
		if err != nil {
			logs.Errorf("Quote type:%v, value:%v, error marshalling array: %v", reflect.TypeOf(v), v, err)
			return strs.Format(fmt, `[]`)
		}
		return strs.Format(fmt, string(bt))
	case []uint8:
		return strs.Format(fmt, string(v))
	case nil:
		return strs.Format(`%s`, "NULL")
	default:
		logs.Errorf("Quote type:%v, value:%v", reflect.TypeOf(v), v)
		return val
	}
}

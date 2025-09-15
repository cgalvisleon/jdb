package jdb

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cgalvisleon/et/console"
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
	quotedChar = fmt.Sprintf(`%s`, char)
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
		return fmt.Sprintf(`%s`, v.ToString())
	case map[string]interface{}:
		return fmt.Sprintf(`%s`, et.Json(v).ToString())
	case time.Time:
		return fmt.Sprintf(`%s`, v.Format("2006-01-02 15:04:05"))
	case []string:
		var r string
		for i, _v := range v {
			if i == 0 {
				r = fmt.Sprintf(`%s`, unquote(_v))
			} else {
				r = fmt.Sprintf(`%s, %s`, r, unquote(_v))
			}
		}
		return fmt.Sprintf(`[%s]`, unquote(r))
	case []interface{}:
		var r string
		for i, _v := range v {
			q := Unquote(_v)
			if i == 0 {
				r = fmt.Sprintf(`%v`, q)
			} else {
				r = fmt.Sprintf(`%s, %v`, r, q)
			}
		}
		return fmt.Sprintf(`[%s]`, r)
	case []uint8:
		return fmt.Sprintf(`%s`, string(v))
	case nil:
		return fmt.Sprintf(`%s`, "NULL")
	default:
		logs.Errorf("Unquoted", "type:%v value:%v", reflect.TypeOf(v), v)
		return val
	}
}

/**
* Quote
* @param val interface{}
* @return any
**/
func Quote(val interface{}) any {
	fm := `'%s'`
	if quotedChar == `"` {
		fm = `"%s"`
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
		return fmt.Sprintf(fm, v.Format("2006-01-02 15:04:05"))
	case []string:
		result := ""
		for _, s := range v {
			result = strs.Append(result, fmt.Sprintf(fm, s), ",")
		}
		return result
	case et.Json:
		return fmt.Sprintf(fm, v.ToString())
	case map[string]interface{}:
		return fmt.Sprintf(fm, et.Json(v).ToString())
	case []et.Json, []interface{}, []map[string]interface{}:
		bt, err := json.Marshal(v)
		if err != nil {
			logs.Errorf("Quote", "type:%v, value:%v, error marshalling array: %v", reflect.TypeOf(v), v, err)
			return fmt.Sprintf(fm, `[]`)
		}
		return fmt.Sprintf(fm, string(bt))
	case []uint8:
		return fmt.Sprintf(fm, string(v))
	case nil:
		return fmt.Sprintf(`%s`, "NULL")
	default:
		logs.Errorf("Quote", "type:%v value:%v", reflect.TypeOf(v), v)
		return val
	}
}

/**
* JsonQuote return a json quote string
* @param val interface{}
* @return interface{}
**/
func JsonQuote(val interface{}) interface{} {
	fm := `'%v'`
	switch v := val.(type) {
	case string:
		v = fmt.Sprintf(`"%s"`, v)
		return fmt.Sprintf(fm, v)
	case int:
		return fmt.Sprintf(fm, v)
	case float64:
		return fmt.Sprintf(fm, v)
	case float32:
		return fmt.Sprintf(fm, v)
	case int16:
		return fmt.Sprintf(fm, v)
	case int32:
		return fmt.Sprintf(fm, v)
	case int64:
		return fmt.Sprintf(fm, v)
	case bool:
		return fmt.Sprintf(fm, v)
	case time.Time:
		return fmt.Sprintf(fm, v.Format("2006-01-02 15:04:05"))
	case et.Json:
		return fmt.Sprintf(fm, v.ToString())
	case map[string]interface{}:
		return fmt.Sprintf(fm, et.Json(v).ToString())
	case []string:
		var r string
		for _, s := range v {
			r = strs.Append(r, fmt.Sprintf(`"%s"`, s), ", ")
		}
		r = fmt.Sprintf(`[%s]`, r)
		return fmt.Sprintf(fm, r)
	case []et.Json, []interface{}, []map[string]interface{}:
		bt, err := json.Marshal(v)
		if err != nil {
			logs.Errorf("JsonQuote", "type:%v, value:%v, error marshalling array: %v", reflect.TypeOf(v), v, err)
			return fmt.Sprintf(fm, `[]`)
		}
		return fmt.Sprintf(fm, string(bt))
	case []uint8:
		return fmt.Sprintf(fm, string(v))
	case nil:
		return fmt.Sprintf(`%s`, "NULL")
	default:
		console.Errorf("JsonQuote type:%v value:%v", reflect.TypeOf(v), v)
		return val
	}
}

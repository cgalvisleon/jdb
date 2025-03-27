package oci

import (
	"reflect"
	"time"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

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
	case []uint8:
		return strs.Format(fmt, string(v))
	case nil:
		return strs.Format(`%s`, "NULL")
	default:
		console.Errorf("Not quoted type:%v value:%v", reflect.TypeOf(v), v)
		return val
	}
}

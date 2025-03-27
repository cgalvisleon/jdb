package cache

import (
	"time"
)

/**
* Set - Set key value
* @params key string, value string, expiration time.Duration
* @return error
**/
func Set(key string, value string, expiration time.Duration) string {
	return conn.Set(key, value, expiration)
}

/**
* Get - Get key value
* @params key, def string
* @return et.KeyValue
**/
func Get(key, def string) string {
	result, exist := conn.Get(key, def)
	if !exist {
		return def
	}

	return result
}

/**
* DeleteCache - Delete key value
* @params key string
* @return error
**/
func DeleteCache(key string) {
	conn.Delete(key)
}

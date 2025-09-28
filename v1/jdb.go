package jdb

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* ConnectTo
* @param name, driver string, userCore bool, params et.Json
* @return (*jdb.Database, error)
**/
func ConnectTo(name, driver string, userCore bool, params et.Json) (*jdb.Database, error) {
	return jdb.ConnectTo(name, driver, userCore, params)
}

/**
* LoadTo
* @param name string
* @return (*jdb.Database, error)
**/
func LoadTo(name string) (*jdb.Database, error) {
	return jdb.LoadTo(name)
}

/**
* Load
* @return (*jdb.Database, error)
**/
func Load() (*jdb.Database, error) {
	return jdb.Load()
}

/**
* Define
* @param definition et.Json
* @return (*jdb.Model, error)
**/
func Define(definition et.Json) (*jdb.Model, error) {
	return jdb.Define(definition)
}

/**
* Select
* @param query et.Json
* @return (*jdb.Ql, error)
**/
func Select(query et.Json) (*jdb.Ql, error) {
	return jdb.Select(query)
}

/**
* Insert
* @param param et.Json
* @return (*jdb.Cmd, error)
**/
func Insert(param et.Json) (*jdb.Cmd, error) {
	return jdb.Insert(param)
}

/**
* Update
* @param param et.Json
* @return (*jdb.Cmd, error)
**/
func Update(param et.Json) (*jdb.Cmd, error) {
	return jdb.Update(param)
}

/**
* Delete
* @param param et.Json
* @return (*jdb.Cmd, error)
**/
func Delete(param et.Json) (*jdb.Cmd, error) {
	return jdb.Delete(param)
}

/**
* Upsert
* @param param et.Json
* @return (*jdb.Cmd, error)
**/
func Upsert(param et.Json) (*jdb.Cmd, error) {
	return jdb.Upsert(param)
}

/**
* Eq
* @param field string, value interface{}
* @return jdb.Condition
**/
func Eq(field string, value interface{}) jdb.Condition {
	return jdb.Eq(field, value)
}

/**
* Neg
* @param field string, value interface{}
* @return jdb.Condition
**/
func Neg(field string, value interface{}) jdb.Condition {
	return jdb.Neg(field, value)
}

/**
* Less
* @param field string, value interface{}
* @return jdb.Condition
**/
func Less(field string, value interface{}) jdb.Condition {
	return jdb.Less(field, value)
}

/**
* LessEq
* @param field string, value interface{}
* @return jdb.Condition
**/
func LessEq(field string, value interface{}) jdb.Condition {
	return jdb.LessEq(field, value)
}

/**
* More
* @param field string, value interface{}
* @return jdb.Condition
**/
func More(field string, value interface{}) jdb.Condition {
	return jdb.More(field, value)
}

/**
* MoreEq
* @param field string, value interface{}
* @return jdb.Condition
**/
func MoreEq(field string, value interface{}) jdb.Condition {
	return jdb.MoreEq(field, value)
}

/**
* Like
* @param field string, value interface{}
* @return jdb.Condition
**/
func Like(field string, value interface{}) jdb.Condition {
	return jdb.Like(field, value)
}

/**
* Ilike
* @param field string, value interface{}
* @return jdb.Condition
**/
func Ilike(field string, value interface{}) jdb.Condition {
	return jdb.Ilike(field, value)
}

/**
* In
* @param field string, value interface{}
* @return jdb.Condition
**/
func In(field string, value interface{}) jdb.Condition {
	return jdb.In(field, value)
}

/**
* NotIn
* @param field string, value interface{}
* @return jdb.Condition
**/
func NotIn(field string, value interface{}) jdb.Condition {
	return jdb.NotIn(field, value)
}

/**
* Is
* @param field string, value interface{}
* @return jdb.Condition
**/
func Is(field string, value interface{}) jdb.Condition {
	return jdb.Is(field, value)
}

/**
* IsNot
* @param field string, value interface{}
* @return jdb.Condition
**/
func IsNot(field string, value interface{}) jdb.Condition {
	return jdb.IsNot(field, value)
}

/**
* Null
* @param field string
* @return jdb.Condition
**/
func Null(field string) jdb.Condition {
	return jdb.Null(field)
}

/**
* NotNull
* @param field string
* @return jdb.Condition
**/
func NotNull(field string) jdb.Condition {
	return jdb.NotNull(field)
}

/**
* Between
* @param field string, value []interface{}
* @return jdb.Condition
**/
func Between(field string, value []interface{}) jdb.Condition {
	return jdb.Between(field, value)
}

/**
* NotBetween
* @param field string, value []interface{}
* @return jdb.Condition
**/
func NotBetween(field string, value []interface{}) jdb.Condition {
	return jdb.NotBetween(field, value)
}

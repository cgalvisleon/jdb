package jdb

import "fmt"

const (
	MSG_DATABASE_REQUIRED   = "database is required"
	MSG_SCHEMA_REQUIRED     = "schema is required"
	MSG_NAME_REQUIRED       = "name is required"
	MSG_TYPE_REQUIRED       = "type is required"
	MSG_DRIVER_REQUIRED     = "driver is required"
	MSG_DRIVER_NOT_FOUND    = "driver %s not found"
	MSG_DATABASE_NOT_FOUND  = "database %s not found"
	MSG_DATA_REQUIRED       = "data is required"
	MSG_FROM_REQUIRED       = "from is required"
	MSG_MODEL_NOT_FOUND     = "model not found"
	MSG_SERIES_NOT_DEFINED  = "series not defined"
	MSG_RECORDS_NOT_DEFINED = "records not defined"
	MSG_RECORD_NOT_FOUND    = "record %s not found"
	MSG_FIELD_NOT_FOUND     = "field not found"
	MSG_FIELD_REQUIRED      = "field %s is required"
	MSG_ATRIB_REQUIRED      = "atrib is required - %s"
)

var (
	ErrModelNotFound = fmt.Errorf(MSG_MODEL_NOT_FOUND)
)

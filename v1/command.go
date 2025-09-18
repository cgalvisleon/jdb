package jdb

import "github.com/cgalvisleon/et/et"

const (
	TypeCommandInsert = "insert"
	TypeCommandUpdate = "update"
	TypeCommandDelete = "delete"
)

type Command struct {
	Before et.Json
	After  et.Json
}

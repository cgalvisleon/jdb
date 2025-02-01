package jdb

type Relation struct {
	Type   TypeColumn
	Column *Column
	To     *Column
}

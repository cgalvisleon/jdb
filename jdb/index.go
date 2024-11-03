package jdb

type Index struct {
	Column *Column
	Unique bool
	Sorted bool
}

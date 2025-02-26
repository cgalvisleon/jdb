package jdb

import "github.com/cgalvisleon/et/et"

type Index struct {
	Column *Column
	Sorted bool
}

/**
* Describe
* @return et.Json
**/
func (i *Index) Describe() et.Json {
	return et.Json{
		"column": i.Column.Name,
		"sorted": i.Sorted,
	}
}

/**
* NewIndex
* @param col *Column
* @param sorted bool
* @return *Index
**/
func NewIndex(col *Column, sorted bool) *Index {
	return &Index{
		Column: col,
		Sorted: sorted,
	}
}

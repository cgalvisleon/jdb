package jdb

type Index struct {
	Column *Column
	Sorted bool
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

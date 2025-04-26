package jdb

type Index struct {
	Column *Column
	Sorted bool
}

/**
* newIndex
* @param col *Column, sorted bool
* @return *Index
**/
func newIndex(col *Column, sorted bool) *Index {
	return &Index{
		Column: col,
		Sorted: sorted,
	}
}

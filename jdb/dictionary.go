package jdb

type Dictionary struct {
	Key     string
	Value   interface{}
	Columns map[string]*Column
}

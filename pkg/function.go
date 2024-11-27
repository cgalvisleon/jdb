package jdb

import "github.com/cgalvisleon/et/strs"

type TypeFunction int

const (
	TpSqlFunction TypeFunction = iota
	TpJsFunction
	TpGoFunction
	TpPythonFunction
)

type Function struct {
	Key          string
	Name         string
	Description  string
	TypeFunction TypeFunction
	Args         []interface{}
	Definition   string
}

/**
* NewFunction
* @param name string
* @param tp TypeFunction
* @return *Function
**/
func NewFunction(name string, tp TypeFunction) *Function {
	return &Function{
		Key:          strs.Uppcase(name),
		Name:         name,
		TypeFunction: tp,
		Args:         make([]interface{}, 0),
		Definition:   ``,
	}
}

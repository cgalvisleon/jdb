package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

type TypeFunction int

const (
	TpSqlFunction TypeFunction = iota
	TpJsFunction
	TpGoFunction
	TpPythonFunction
)

type Function struct {
	Id           string
	Name         string
	Description  string
	TypeFunction TypeFunction
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
		Id:           utility.UUID(),
		Name:         name,
		TypeFunction: tp,
		Definition:   ``,
	}
}

/**
* Exec
* @param params et.Json
* @return et.Json
**/
func (s *Function) Exec(params et.Json) et.Json {
	return et.Json{}
}

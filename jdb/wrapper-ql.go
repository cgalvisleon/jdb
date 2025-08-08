package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/dop251/goja"
)

type QlWrapper struct {
	vm      *goja.Runtime
	db      *DB
	ql      *Ql
	command *Command
	tx      *Tx
}

/**
* Ql
* @param call goja.FunctionCall
* @return goja.Value
**/
func (s *QlWrapper) Ql(call goja.FunctionCall) goja.Value {
	arg := call.Argument(0).Export()
	where, ok := arg.(map[string]interface{})
	if !ok {
		panic("query is required")
	}

	if s.ql == nil {
		s.ql = NewQl(s.db)
	}

	json, err := s.ql.QueryTx(s.tx, where)
	if err != nil {
		panic(err)
	}

	return s.vm.ToValue(json)
}

/**
* Insert
* @param call goja.FunctionCall
* @return goja.Value
**/
func (s *QlWrapper) Insert(call goja.FunctionCall) goja.Value {
	arg := call.Argument(0).Export()
	query, ok := arg.(map[string]interface{})
	if !ok {
		panic("data is required")
	}

	var from string
	var data et.Json
	var where et.Json
	var returns et.Json
	var isDebug bool
	for key, val := range query {
		if key == "where" {
			where = val.(et.Json)
		} else if key == "returns" {
			returns = val.(et.Json)
		} else if key == "is_debug" {
			isDebug = val.(bool)
		} else {
			from = key
			data = val.(et.Json)
		}
	}

	model := s.db.GetModel(from)
	s.command = model.
		Insert(data).
		setWheres(where).
		setReturns(returns).
		setIsDebug(isDebug)

	result, err := s.command.
		Exec()
	if err != nil {
		panic(err)
	}

	return s.vm.ToValue(result.ToMap())
}

/**
* NewQlWrapper
* @param vm *goja.Runtime, db *DB
* @return map[string]interface{}
**/
func NewQlWrapper(vm *goja.Runtime, db *DB) map[string]interface{} {
	wrapper := &QlWrapper{
		vm: vm,
		db: db,
	}

	return map[string]interface{}{
		"Ql":     wrapper.Ql,
		"Insert": wrapper.Insert,
	}
}

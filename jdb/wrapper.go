package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
	"github.com/dop251/goja"
)

type ModelWrapper struct {
	Db         *DB
	ProjectId  string
	InstanceId string
	ContextId  string
	Model      *Model
	Ql         *Ql
	tx         *Tx
	where      et.Json
	command    *Command
	isDebug    bool
}

/**
* NewModelWrapper
* @param vm *goja.Runtime, m *ModelWrapper
* @return map[string]interface{}
**/
func NewModelWrapper(vm *goja.Runtime, m *ModelWrapper) map[string]interface{} {
	wrapper := map[string]interface{}{}

	wrapper["select"] = func(call goja.FunctionCall) goja.Value {
		arr := call.Argument(0).Export()
		stringRaw, ok := arr.([]interface{})
		if !ok {
			panic(vm.NewTypeError("select espera un array"))
		}
		fields := make([]interface{}, len(stringRaw))
		copy(fields, stringRaw)
		m.Ql.Select(fields...)
		return vm.ToValue(wrapper)
	}

	wrapper["data"] = func(call goja.FunctionCall) goja.Value {
		arr := call.Argument(0).Export()
		stringRaw, ok := arr.([]interface{})
		if !ok {
			panic(vm.NewTypeError("select espera un array"))
		}
		fields := make([]interface{}, len(stringRaw))
		copy(fields, stringRaw)
		m.Ql.Data(fields...)
		return vm.ToValue(wrapper)
	}

	wrapper["from"] = func(call goja.FunctionCall) goja.Value {
		name := call.Argument(0).String()
		if !utility.ValidStr(name, 3, []string{}) {
			panic(vm.NewTypeError("from espera un string"))
		}
		model, err := LoadModel(m.Db, name)
		if err != nil {
			panic(vm.NewTypeError("from espera un string"))
		}
		m.Ql.From(model.Name)
		return vm.ToValue(wrapper)
	}

	wrapper["join"] = func(call goja.FunctionCall) goja.Value {
		arr := call.Argument(0).Export()
		filtros, ok := arr.([]map[string]interface{})
		if !ok {
			panic(vm.NewTypeError("where espera un objeto"))
		}
		m.Ql.Join(filtros)
		return vm.ToValue(wrapper)
	}

	wrapper["where"] = func(call goja.FunctionCall) goja.Value {
		arr := call.Argument(0).Export()
		filtros, ok := arr.(map[string]interface{})
		if !ok {
			panic(vm.NewTypeError("where espera un objeto"))
		}
		m.where = filtros
		return vm.ToValue(wrapper)
	}

	wrapper["groupBy"] = func(call goja.FunctionCall) goja.Value {
		arr := call.Argument(0).Export()
		stringRaw, ok := arr.([]interface{})
		if !ok {
			panic(vm.NewTypeError("select espera un array"))
		}
		fields := make([]string, len(stringRaw))
		for i, v := range stringRaw {
			fields[i] = fmt.Sprint(v)
		}
		m.Ql.GroupBy(fields...)
		return vm.ToValue(wrapper)
	}

	wrapper["having"] = func(call goja.FunctionCall) goja.Value {
		arr := call.Argument(0).Export()
		filtros, ok := arr.(map[string]interface{})
		if !ok {
			panic(vm.NewTypeError("where espera un objeto"))
		}
		m.Ql.SetHavings(filtros)
		return vm.ToValue(wrapper)
	}

	wrapper["orderBy"] = func(call goja.FunctionCall) goja.Value {
		arr := call.Argument(0).Export()
		filtros, ok := arr.(map[string]interface{})
		if !ok {
			panic(vm.NewTypeError("where espera un objeto"))
		}
		m.Ql.SetOrderBy(filtros)
		return vm.ToValue(wrapper)
	}

	wrapper["debug"] = func(call goja.FunctionCall) goja.Value {
		m.isDebug = true
		return vm.ToValue(wrapper)
	}

	wrapper["page"] = func(call goja.FunctionCall) goja.Value {
		page := call.Argument(0).ToInteger()
		if page == 0 {
			page = 1
		}
		p := int(page)
		m.Ql.SetPage(p)
		return vm.ToValue(wrapper)
	}

	wrapper["limit"] = func(call goja.FunctionCall) goja.Value {
		limit := call.Argument(0).ToInteger()
		if limit == 0 {
			limit = 1
		}
		l := int(limit)
		m.Ql.SetWheres(m.where)
		m.Ql.SetDebug(m.isDebug)
		result, err := m.Ql.SetLimitTx(m.tx, l)
		if err != nil {
			panic(vm.NewTypeError(err.Error()))
		}

		return vm.ToValue(result.ToMap())
	}

	wrapper["list"] = func(call goja.FunctionCall) goja.Value {
		arr := call.Argument(0).Export()
		intRaw, ok := arr.([]interface{})
		if !ok {
			panic(vm.NewTypeError("list espera un array de enteros"))
		}
		if len(intRaw) == 0 {
			panic(vm.NewTypeError("list espera un array de enteros"))
		}

		page := int(intRaw[0].(int))
		rows := int(intRaw[1].(int))
		m.Ql.SetWheres(m.where)
		m.Ql.SetDebug(m.isDebug)
		result, err := m.Ql.ListTx(m.tx, page, rows)
		if err != nil {
			panic(vm.NewTypeError(err.Error()))
		}

		return vm.ToValue(result.ToMap())
	}

	wrapper["first"] = func(call goja.FunctionCall) goja.Value {
		limit := call.Argument(0).ToInteger()
		if limit == 0 {
			limit = 1
		}
		l := int(limit)
		m.Ql.SetWheres(m.where)
		m.Ql.SetDebug(m.isDebug)
		result, err := m.Ql.FirstTx(m.tx, l)
		if err != nil {
			panic(vm.NewTypeError(err.Error()))
		}

		return vm.ToValue(result.ToMap())
	}

	wrapper["all"] = func(call goja.FunctionCall) goja.Value {
		m.Ql.SetWheres(m.where)
		m.Ql.SetDebug(m.isDebug)
		result, err := m.Ql.AllTx(m.tx)
		if err != nil {
			panic(vm.NewTypeError(err.Error()))
		}

		return vm.ToValue(result.ToMap())
	}

	wrapper["insert"] = func(call goja.FunctionCall) goja.Value {
		arr := call.Argument(0).Export()
		data, ok := arr.(map[string]interface{})
		if !ok {
			panic(vm.NewTypeError("insert espera un objeto"))
		}
		m.command = m.Model.Insert(data)
		return vm.ToValue(wrapper)
	}

	wrapper["update"] = func(call goja.FunctionCall) goja.Value {
		arr := call.Argument(0).Export()
		data, ok := arr.(map[string]interface{})
		if !ok {
			panic(vm.NewTypeError("update espera un objeto"))
		}
		m.command = m.Model.Update(data)
		return vm.ToValue(wrapper)
	}

	wrapper["delete"] = func(call goja.FunctionCall) goja.Value {
		arr := call.Argument(0).Export()
		filtros, ok := arr.(map[string]interface{})
		if !ok {
			panic(vm.NewTypeError("delete espera un filtro"))
		}
		m.command = m.Model.Delete("")
		m.where = filtros
		return vm.ToValue(wrapper)
	}

	wrapper["exec"] = func(call goja.FunctionCall) goja.Value {
		if m.command == nil {
			panic(vm.NewTypeError("exec espera un comando"))
		}
		if m.tx == nil {
			m.tx = NewTx()
		}
		m.Ql.SetWheres(m.where)
		m.Ql.SetDebug(m.isDebug)
		result, err := m.command.ExecTx(m.tx)
		if err != nil {
			panic(vm.NewTypeError(err.Error()))
		}
		return vm.ToValue(result.ToMap())
	}

	wrapper["execOne"] = func(call goja.FunctionCall) goja.Value {
		if m.command == nil {
			panic(vm.NewTypeError("exec espera un comando"))
		}
		if m.tx == nil {
			m.tx = NewTx()
		}
		m.Ql.SetWheres(m.where)
		m.Ql.SetDebug(m.isDebug)
		result, err := m.command.OneTx(m.tx)
		if err != nil {
			panic(vm.NewTypeError(err.Error()))
		}
		return vm.ToValue(result.ToMap())
	}

	return wrapper
}

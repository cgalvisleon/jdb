package jdb

import (
	"slices"
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/reg"
)

type TypeCommand int

const (
	Insert TypeCommand = iota
	Update
	Delete
	Upsert
	Delsert
	Sync
)

func (s TypeCommand) Str() string {
	switch s {
	case Insert:
		return "insert"
	case Update:
		return "update"
	case Delete:
		return "delete"
	case Upsert:
		return "upsert"
	case Delsert:
		return "delsert"
	case Sync:
		return "sync"
	default:
		return "No command"
	}
}

type Function func() error
type DataFunction func(data et.Json) error
type DataFunctionTx func(tx *Tx, data et.Json) error

type Command struct {
	*QlWhere
	Id           string              `json:"id"`
	tx           *Tx                 `json:"-"`
	Command      TypeCommand         `json:"command"`
	Db           *DB                 `json:"-"`
	From         *Model              `json:"-"`
	Data         []et.Json           `json:"data"`
	Values       []map[string]*Field `json:"values"`
	Returns      []*Field            `json:"returns"`
	Sql          string              `json:"sql"`
	Result       et.Items            `json:"result"`
	beforeInsert []DataFunctionTx    `json:"-"`
	beforeUpdate []DataFunctionTx    `json:"-"`
	beforeDelete []Function          `json:"-"`
	afterInsert  []DataFunctionTx    `json:"-"`
	afterUpdate  []DataFunctionTx    `json:"-"`
	afterDelete  []Function          `json:"-"`
	isSync       bool                `json:"-"`
}

/**
* NewCommand
* @param model *Model, data []et.Json, command TypeCommand
* @return *Command
**/
func NewCommand(model *Model, data []et.Json, command TypeCommand) *Command {
	result := &Command{
		Id:           reg.GenId("command"),
		Command:      command,
		Db:           model.Db,
		From:         model,
		QlWhere:      NewQlWhere(),
		Data:         data,
		Values:       make([]map[string]*Field, 0),
		beforeInsert: []DataFunctionTx{},
		beforeUpdate: []DataFunctionTx{},
		beforeDelete: []Function{},
		afterInsert:  []DataFunctionTx{},
		afterUpdate:  []DataFunctionTx{},
		afterDelete:  []Function{},
		Returns:      []*Field{},
		Result:       et.Items{},
	}
	result.beforeInsert = append(result.beforeInsert, result.beforeInsertDefault)
	result.beforeUpdate = append(result.beforeUpdate, result.beforeUpdateDefault)

	return result
}

/**
* setTx
* @param tx *Tx
* @return *Command
**/
func (s *Command) setTx(tx *Tx) *Command {
	s.tx = tx

	return s
}

/**
* setSync
* @param isSync bool
* @return *Command
**/
func (s *Command) setSync() *Command {
	s.isSync = true
	return s
}

/**
* Describe
* @return et.Json
**/
func (s *Command) Describe() et.Json {
	result, err := et.Object(s)
	if err != nil {
		return et.Json{}
	}

	values := []et.Json{}
	for _, val := range s.Values {
		for _, field := range val {
			values = append(values, field.describe())
		}
	}

	result["command"] = s.Command.Str()
	result["wheres"] = s.listWheres()
	result["values"] = values

	return result
}

/**
* getField
* @param name string, isCreate bool
* @return *Field
**/
func (s *Command) getField(name string, isCreate bool) *Field {
	return s.From.getField(name, isCreate)
}

/**
* validator
* validate this val is a field or basic type
* @return interface{}
**/
func (s *Command) validator(val interface{}) interface{} {
	switch v := val.(type) {
	case string:
		if strings.HasPrefix(v, ":") {
			v = strings.TrimPrefix(v, ":")
			field := s.getField(v, false)
			if field != nil {
				return field
			}
			return nil
		}

		if strings.HasPrefix(v, "$") {
			v = strings.TrimPrefix(v, "$")
			return v
		}

		v = strings.Replace(v, `\\:`, `\:`, 1)
		v = strings.Replace(v, `\:`, `:`, 1)
		v = strings.Replace(v, `\\$`, `\$`, 1)
		v = strings.Replace(v, `\$`, `$`, 1)
		field := s.getField(v, false)
		if field != nil {
			return field
		}

		return v
	default:
		if v == nil {
			return "nil"
		}

		return v
	}
}

/**
* setWhere
* @param setWheres et.Json
* @return *Command
**/
func (s *Command) setWheres(wheres et.Json) *Command {
	and := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.and(key)
				s.setValue(val.Json(key), s.validator)
			}
		}
	}

	or := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.or(key)
				s.setValue(val.Json(key), s.validator)
			}
		}
	}

	for key := range wheres {
		if slices.Contains([]string{"and", "AND", "or", "OR"}, key) {
			continue
		}

		s.Where(key).setValue(wheres.Json(key), s.validator)
	}

	for key := range wheres {
		switch key {
		case "and", "AND":
			vals := wheres.ArrayJson(key)
			and(vals)
		case "or", "OR":
			vals := wheres.ArrayJson(key)
			or(vals)
		}
	}

	return s
}

/**
* listWheres
* @return et.Json
**/
func (s *Command) listWheres() et.Json {
	return s.QlWhere.listWheres()
}

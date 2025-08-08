package jdb

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
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
type DataFunction func(data et.Json)
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
	Current      et.Items            `json:"current"`
	CurrentMap   map[string]et.Json  `json:"current_map"`
	ResultMap    map[string]et.Json  `json:"result_map"`
	beforeInsert []DataFunctionTx    `json:"-"`
	beforeUpdate []DataFunctionTx    `json:"-"`
	beforeDelete []DataFunctionTx    `json:"-"`
	afterInsert  []DataFunctionTx    `json:"-"`
	afterUpdate  []DataFunctionTx    `json:"-"`
	afterDelete  []DataFunctionTx    `json:"-"`
	isSync       bool                `json:"-"`
}

/**
* NewCommand
* @param model *Model, data []et.Json, command TypeCommand
* @return *Command
**/
func NewCommand(model *Model, data []et.Json, command TypeCommand) *Command {
	result := &Command{
		Id:           utility.UUID(),
		Command:      command,
		Db:           model.Db,
		From:         model,
		Data:         data,
		Values:       make([]map[string]*Field, 0),
		beforeInsert: []DataFunctionTx{},
		beforeUpdate: []DataFunctionTx{},
		beforeDelete: []DataFunctionTx{},
		afterInsert:  []DataFunctionTx{},
		afterUpdate:  []DataFunctionTx{},
		afterDelete:  []DataFunctionTx{},
		Returns:      []*Field{},
		Result:       et.Items{},
		Current:      et.Items{},
		CurrentMap:   make(map[string]et.Json),
		ResultMap:    make(map[string]et.Json),
	}
	result.QlWhere = newQlWhere(result.validator)
	result.IsDebug = model.IsDebug
	result.beforeInsert = append(result.beforeInsert, result.beforeInsertDefault)
	result.beforeUpdate = append(result.beforeUpdate, result.beforeUpdateDefault)
	result.beforeDelete = append(result.beforeDelete, result.beforeDeleteDefault)

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
* serialize
* @return []byte, error
**/
func (s *Command) serialize() ([]byte, error) {
	result, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

/**
* Describe
* @return et.Json
**/
func (s *Command) Describe() et.Json {
	definition, err := s.serialize()
	if err != nil {
		return et.Json{}
	}

	result := et.Json{}
	err = json.Unmarshal(definition, &result)
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
	result["wheres"] = s.getWheres()
	result["values"] = values

	return result
}

/**
* getField
* @param name string
* @return *Field
**/
func (s *Command) getField(name string) *Field {
	return s.From.getField(name, false)
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
			field := s.getField(v)
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
		field := s.getField(v)
		if field != nil {
			return field
		}

		return v
	case *Field:
		return v
	case Field:
		return v
	case *Column:
		return v.GetField()
	case Column:
		return v.GetField()
	case []interface{}:
		return v
	case []string:
		return v
	case []et.Json:
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
				s.And(key).setValue(val.Json(key))
			}
		}
	}

	or := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.Or(key).setValue(val.Json(key))
			}
		}
	}

	for key := range wheres {
		key = strings.ToLower(key)
		if slices.Contains([]string{"and", "or"}, key) {
			continue
		}

		val := wheres.Json(key)
		s.Where(key).setValue(val)
	}

	for key := range wheres {
		switch strings.ToLower(key) {
		case "and":
			vals := wheres.ArrayJson(key)
			and(vals)
		case "or":
			vals := wheres.ArrayJson(key)
			or(vals)
		}
	}

	return s
}

/**
* setReturns
* @param returns et.Json
* @return *Command
**/
func (s *Command) setReturns(returns et.Json) *Command {
	for key := range returns {
		s.Returns = append(s.Returns, s.From.getField(key, false))
	}

	return s
}

/**
* setIsDebug
* @param isDebug bool
* @return *Command
**/
func (s *Command) setIsDebug(isDebug bool) *Command {
	s.IsDebug = isDebug

	return s
}

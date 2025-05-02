package jdb

import (
	"slices"
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

type TypeCommand int

const (
	Insert TypeCommand = iota
	Update
	Upsert
	Delete
	Bulk
	Sync
)

type DataFunction func(data et.Json) et.Json

type Command struct {
	*QlWhere
	tx           *Tx                 `json:"-"`
	Command      TypeCommand         `json:"command"`
	Db           *DB                 `json:"-"`
	From         *Model              `json:"-"`
	Data         []et.Json           `json:"data"`
	Values       []map[string]*Field `json:"values"`
	Returns      []*Field            `json:"returns"`
	Sql          string              `json:"sql"`
	Result       et.Items            `json:"result"`
	beforeInsert []DataFunction      `json:"-"`
	beforeUpdate []DataFunction      `json:"-"`
	isUndo       bool                `json:"-"`
	isSync       bool                `json:"-"`
}

/**
* NewCommand
* @param model *Model, data []et.Json, command TypeCommand
* @return *Command
**/
func NewCommand(model *Model, data []et.Json, command TypeCommand) *Command {
	result := &Command{
		Command:      command,
		Db:           model.Db,
		From:         model,
		QlWhere:      NewQlWhere(),
		Data:         data,
		Values:       []map[string]*Field{},
		beforeInsert: []DataFunction{},
		beforeUpdate: []DataFunction{},
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

	return result
}

/**
* Tx
* @return *Tx
**/
func (s *Command) Tx() *Tx {
	return s.tx
}

/**
* ExecTx
* @param tx *Tx
* @return et.Items, error
**/
func (s *Command) ExecTx(tx *Tx) (et.Items, error) {
	s.setTx(tx)

	switch s.Command {
	case Insert:
		if len(s.Data) == 0 {
			return et.Items{}, mistake.New(MSG_NOT_DATA)
		}

		err := s.inserted()
		if err != nil {
			return et.Items{}, err
		}
	case Update:
		if len(s.Data) == 0 {
			return et.Items{}, mistake.New(MSG_NOT_DATA)
		}

		err := s.updated()
		if err != nil {
			return et.Items{}, err
		}
	case Upsert:
		err := s.upsert()
		if err != nil {
			return et.Items{}, err
		}
	case Delete:
		err := s.deleted()
		if err != nil {
			return et.Items{}, err
		}
	case Bulk:
		if len(s.Data) == 0 {
			return et.Items{}, mistake.New(MSG_NOT_DATA)
		}

		err := s.bulk()
		if err != nil {
			return et.Items{}, err
		}
	case Sync:
		if len(s.Data) == 0 {
			return et.Items{}, mistake.New(MSG_NOT_DATA)
		}

		err := s.sync()
		if err != nil {
			return et.Items{}, err
		}
	default:
		return et.Items{}, mistake.New(MSG_NOT_COMMAND)
	}

	return s.Result, nil
}

/**
* OneTx
* @param tx *Tx
* @return et.Item, error
**/
func (s *Command) OneTx(tx *Tx) (et.Item, error) {
	result, err := s.ExecTx(tx)
	if err != nil {
		return et.Item{}, err
	}

	return result.First(), nil
}

/**
* Exec
* @return et.Items, error
**/
func (s *Command) Exec() (et.Items, error) {
	return s.ExecTx(nil)
}

/**
* One
* @return et.Item, error
**/
func (s *Command) One() (et.Item, error) {
	return s.OneTx(nil)
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
		if strings.HasPrefix(v, "$") {
			v = strings.TrimPrefix(v, "$")
			field := s.getField(v, false)
			if field != nil {
				return field
			}
			return nil
		}
		result := s.getField(v, false)
		if result != nil {
			return result
		}

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

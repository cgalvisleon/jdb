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
	Delete
	Bulk
	Undo
	Sync
)

type Command struct {
	*QlWhere
	Command  TypeCommand
	Db       *DB
	From     *Model
	Data     []et.Json
	Values   []map[string]*Field
	Returns  []*Field
	Undo     et.Json
	Sql      string
	Result   et.Items
	rollback bool
}

/**
* NewCommand
* @param model *Model
* @param data []et.Json
* @param command TypeCommand
* @return *Command
**/
func NewCommand(model *Model, data []et.Json, command TypeCommand) *Command {
	result := &Command{
		Command: command,
		Db:      model.Db,
		From:    model,
		QlWhere: NewQlWhere(),
		Data:    data,
		Undo:    et.Json{},
		Values:  []map[string]*Field{},
		Returns: []*Field{},
		Result:  et.Items{},
	}

	return result
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
* Exec
* @return et.Items, error
**/
func (s *Command) Exec() (et.Items, error) {
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
	case Delete:
		err := s.delete()
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
	case Undo:
		err := s.undo()
		if err != nil {
			return et.Items{}, err
		}
	default:
		return et.Items{}, mistake.New(MSG_NOT_COMMAND)
	}

	return s.Result, nil
}

/**
* One
* @return et.Item, error
**/
func (s *Command) One() (et.Item, error) {
	result, err := s.Exec()
	if err != nil {
		return et.Item{}, err
	}

	if !result.Ok {
		return et.Item{Result: et.Json{}}, nil
	}

	return et.Item{
		Ok:     true,
		Result: result.Result[0],
	}, nil
}

/**
* Rollback
* @return et.Items, error
**/
func (s *Command) Rollback() error {
	s.rollback = true
	_, err := s.Exec()
	if err != nil {
		return err
	}

	return nil
}

/**
* Return
* @param fields ...string
* @return *Command
**/
func (s *Command) Return(fields ...string) *Command {
	for _, name := range fields {
		field := s.getField(name, true)
		if field == nil {
			continue
		}

		s.Returns = append(s.Returns, field)
	}

	return s
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
* @param wheres et.Json
* @return *Command
**/
func (s *Command) setWhere(wheres et.Json) *Command {
	if len(wheres) == 0 {
		return s
	}

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

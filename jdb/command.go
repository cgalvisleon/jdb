package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
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
* asField
* @param field *Field
* @return string
**/
func (s *Command) asField(field *Field) string {
	if s.From == nil {
		return field.Name
	}

	return strs.Format("%s.%s", field.Table, field.Name)
}

/**
* setWhere
* @param wheres et.Json
* @return *Command
**/
func (s *Command) setWhere(wheres et.Json) *Command {
	s.QlWhere.setWheres(wheres, s.getField)

	return s
}

/**
* listWheres
* @return et.Json
**/
func (s *Command) listWheres() et.Json {
	return s.QlWhere.listWheres(s.asField)
}

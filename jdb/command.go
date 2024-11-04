package jdb

import (
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
)

type TypeCommand int

const (
	Insert TypeCommand = iota
	Update
	Delete
)

type Command struct {
	Db      *Database
	Model   *Model
	Data    []et.Json
	Wheres  []*LinqWhere
	Command TypeCommand
	Show    bool
	Sql     string
	Result  et.Items
}

/**
* NewCommand
* @param model *Model
* @param data []et.Json
* @param command TypeCommand
* @return *Command
**/
func NewCommand(model *Model, data []et.Json, command TypeCommand) *Command {
	return &Command{
		Db:      model.Db,
		Model:   model,
		Data:    data,
		Wheres:  make([]*LinqWhere, 0),
		Command: command,
		Show:    false,
		Sql:     "",
		Result:  et.Items{},
	}
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
* getColumn
* @param col interface{}
* @return *LinqSelect
**/
func (s *Command) getColumn(col interface{}) *LinqSelect {
	switch v := col.(type) {
	case Column:
		from := &LinqFrom{
			Model: *s.Model,
			As:    s.Model.Table,
		}

		return &LinqSelect{
			From:  from,
			Field: v.Field,
		}
	case *Column:
		from := &LinqFrom{
			Model: *s.Model,
			As:    s.Model.Table,
		}

		return &LinqSelect{
			From:  from,
			Field: v.Field,
		}
	case string:
		list := strings.Split(v, ".")
		if len(list[0]) == 0 {
			return nil
		}

		from := &LinqFrom{
			Model: *s.Model,
			As:    s.Model.Table,
		}

		if len(list[1]) == 0 {
			return &LinqSelect{
				From:  from,
				Field: strs.Uppcase(list[0]),
			}
		}

		return nil
	default:
		return nil
	}
}

/**
* Insert
* @param data []et.Json
* @return *Command
**/
func (s *Model) Insert(data []et.Json) *Command {
	return NewCommand(s, data, Insert)
}

/**
* Update
* @param data []et.Json
* @return *Command
**/
func (s *Model) Update(data []et.Json) *Command {
	return NewCommand(s, data, Update)
}

/**
* Delete
* @return *Command
**/
func (s *Model) Delete() *Command {
	return NewCommand(s, []et.Json{}, Delete)
}

/**
* Debug
* @return *Command
**/
func (s *Command) Debug() *Command {
	s.Show = true
	return s
}

/**
* Exec
* @return et.Items, error
**/
func (s *Command) Exec() (et.Items, error) {
	result, err := (*s.Db.Driver).Command(s)
	if s.Show {
		logs.Debug(s.Describe().ToString())
	}
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
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
		return et.Item{}, nil
	}

	return et.Item{
		Ok:     true,
		Result: result.Result[0],
	}, nil
}

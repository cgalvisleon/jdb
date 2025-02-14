package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/utility"
)

type TypeCommand int

const (
	Insert TypeCommand = iota
	Update
	Delete
	Bulk
	Undo
)

type Value struct {
	Columns et.Json
	Atribs  et.Json
	Data    et.Json
}

func NewValue() *Value {
	return &Value{
		Columns: et.Json{},
		Atribs:  et.Json{},
		Data:    et.Json{},
	}
}

type Command struct {
	*QlFilter
	Db         *DB
	TypeSelect TypeSelect
	From       *QlFrom
	Command    TypeCommand
	Origin     []et.Json
	Values     []*Value
	Undo       *UndoRecord
	Sql        string
	Result     et.Items
	history    bool
}

/**
* NewCommand
* @param model *Model
* @param data []et.Json
* @param command TypeCommand
* @return *Command
**/
func NewCommand(model *Model, data []et.Json, command TypeCommand) *Command {
	tp := Select
	if model.SourceField != nil {
		tp = Data
	}
	result := &Command{
		TypeSelect: tp,
		Command:    command,
		Origin:     data,
		Values:     make([]*Value, 0),
		Sql:        "",
		Result:     et.Items{},
	}
	result.QlFilter = &QlFilter{
		main:   result,
		Wheres: make([]*QlWhere, 0),
	}
	result.addFrom(model)

	return result
}

/**
* addFrom
* @param m *Model
* @return *QlFrom
**/
func (s *Command) addFrom(m *Model) *QlFrom {
	s.Db = m.Db
	s.From = &QlFrom{
		Model: m,
		As:    "",
	}

	return s.From
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
	switch s.Command {
	case Insert:
		if len(s.Origin) == 0 {
			return et.Items{}, mistake.New(MSG_NOT_DATA)
		}

		err := s.inserted()
		if err != nil {
			return et.Items{}, err
		}
	case Update:
		if len(s.Origin) == 0 {
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
		if len(s.Origin) == 0 {
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
		return et.Item{}, nil
	}

	return et.Item{
		Ok:     true,
		Result: result.Result[0],
	}, nil
}

/**
* GetReturn
* @param name string
* @return *QlSelect
**/
func (s *Command) GetReturn(name string) *QlSelect {
	field := s.From.GetField(name, true)
	if field == nil {
		return nil
	}

	return NewQlSelect(s.From, field)
}

/**
* Commands
* @param command et.Json
* @return et.Items, error
**/
func Commands(command et.Json) (et.Items, error) {
	if command.IsEmpty() {
		return et.Items{}, mistake.New(MSG_QUERY_EMPTY)
	}

	from := command.Str("from")
	if !utility.ValidStr(from, 0, []string{""}) {
		return et.Items{}, mistake.New(MSG_QUERY_FROM_REQUIRED)
	}

	model := Jdb.Models[from]
	if model == nil {
		return et.Items{}, mistake.Newf(MSG_MODEL_NOT_FOUND, from)
	}

	return model.
		Command(command)
}

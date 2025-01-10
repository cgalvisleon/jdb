package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

type TypeCommand int

const (
	Insert TypeCommand = iota
	Update
	Delete
	Bulk
)

type Command struct {
	*LinqFilter
	Db         *DB
	TypeSelect TypeSelect `json:"type_select"`
	From       *LinqFrom
	Command    TypeCommand
	Origin     []et.Json
	Atribs     et.Json
	Fields     et.Json
	New        *et.Json
	Sql        string
	Result     et.Items
	Show       bool
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
		Atribs:     et.Json{},
		Fields:     et.Json{},
		Origin:     data,
		New:        &et.Json{},
		Show:       false,
		Sql:        "",
		Result:     et.Items{},
	}
	result.LinqFilter = &LinqFilter{
		main:   result,
		Wheres: make([]*LinqWhere, 0),
	}
	result.addFrom(model)

	return result
}

/**
* addFrom
* @param m *Model
* @return *LinqFrom
**/
func (s *Command) addFrom(m *Model) *LinqFrom {
	s.Db = m.Db
	s.From = &LinqFrom{
		Model:   m,
		As:      "",
		Selects: make([]*LinqSelect, 0),
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

func (s *Command) consolidate(data et.Json) et.Json {
	for k, v := range data {
		field := s.From.GetField(k, true)
		if field != nil {
			(*s.New)[k] = v
			switch field.Column.TypeColumn {
			case TpAtribute:
				s.Atribs[k] = v
			case TpColumn:
				s.Fields[k] = v
			}
		} else if !s.From.Integrity {
			(*s.New)[k] = v
			s.Atribs[k] = v
		}
	}

	return (*s.New)
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

		result, err := s.inserted()
		if err != nil {
			return et.Items{}, err
		}
		s.Result.Add(result.Result)
	case Update:
		if len(s.Origin) == 0 {
			return et.Items{}, mistake.New(MSG_NOT_DATA)
		}

		result, err := s.updated()
		if err != nil {
			return et.Items{}, err
		}
		s.Result = result
	case Delete:
		result, err := s.delete()
		if err != nil {
			return et.Items{}, err
		}
		s.Result = result
	case Bulk:
		if len(s.Origin) == 0 {
			return et.Items{}, mistake.New(MSG_NOT_DATA)
		}

		result, err := s.bulk()
		if err != nil {
			return et.Items{}, err
		}
		s.Result.Add(result.Result)
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
* @return *LinqSelect
**/
func (s *Command) GetReturn(name string) *LinqSelect {
	field := s.From.GetField(name, true)
	if field == nil {
		return nil
	}

	return NewLinqSelect(s.From, field)
}

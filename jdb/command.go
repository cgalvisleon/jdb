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
	*Model
	*LinqFilter
	TypeSelect TypeSelect `json:"type_select"`
	Command    TypeCommand
	Origin     []et.Json
	Atribs     et.Json
	Fields     et.Json
	Key        string
	New        *et.Json
	Returns    []*LinqSelect
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
		Model:      model,
		TypeSelect: tp,
		Command:    command,
		Atribs:     et.Json{},
		Fields:     et.Json{},
		Origin:     data,
		New:        &et.Json{},
		Returns:    make([]*LinqSelect, 0),
		Show:       false,
		Sql:        "",
		Result:     et.Items{},
	}
	result.LinqFilter = &LinqFilter{
		main:   result,
		Wheres: make([]*LinqWhere, 0),
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

func (s *Command) consolidate(data et.Json) et.Json {
	for k, v := range data {
		field := s.GetField(k)
		if field != nil {
			(*s.New)[k] = v
			switch field.Column.TypeColumn {
			case TpAtribute:
				s.Atribs[k] = v
			case TpColumn:
				s.Fields[k] = v
			}
		} else if !s.Integrity {
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
			return et.Items{}, mistake.New("Data not found")
		}

		for _, data := range s.Origin {
			result, err := s.inserted(data)
			if err != nil {
				return et.Items{}, err
			}
			s.Result.Result = append(s.Result.Result, result.Result)
			s.Result.Count++
			s.Result.Ok = true
		}
	case Update:
		if len(s.Origin) == 0 {
			return et.Items{}, mistake.New("Data not found")
		}

		current, err := s.Db.Current(s)
		if err != nil {
			return et.Items{}, err
		}
		for _, old := range current.Result {
			s.Key = old.ValStr("-1", SystemKeyField.Str())
			if s.Key == "-1" {
				continue
			}
			result, err := s.updated(old, s.Origin[0])
			if err != nil {
				return et.Items{}, err
			}
			s.Result.Result = append(s.Result.Result, result.Result)
			s.Result.Count++
			s.Result.Ok = true
		}
	case Delete:
		current, err := s.Db.Current(s)
		if err != nil {
			return et.Items{}, err
		}
		for _, old := range current.Result {
			s.Key = old.ValStr("-1", SystemKeyField.Str())
			if s.Key == "-1" {
				continue
			}
			result, err := s.delete(old)
			if err != nil {
				return et.Items{}, err
			}
			s.Result.Result = append(s.Result.Result, result.Result)
			s.Result.Count++
			s.Result.Ok = true
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
* @return *LinqSelect
**/
func (s *Command) GetReturn(name string) *LinqSelect {
	field := s.GetField(name)
	if field == nil {
		return nil
	}

	from := &LinqFrom{
		Model: s.Model,
		As:    s.Model.Table,
	}

	return NewLinqSelect(from, field)
}

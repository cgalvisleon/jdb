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
	Command TypeCommand
	Data    []et.Json
	Atribs  et.Json
	Fields  et.Json
	Key     string
	New     *et.Json
	Wheres  []*LinqWhere
	Returns []*LinqSelect
	Sql     string
	Result  et.Items
	Show    bool
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
		Model:   model,
		Command: command,
		Atribs:  et.Json{},
		Fields:  et.Json{},
		Data:    data,
		New:     &et.Json{},
		Wheres:  make([]*LinqWhere, 0),
		Returns: make([]*LinqSelect, 0),
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

func (s *Command) consolidate(data et.Json) et.Json {
	for k, v := range data {
		col := s.GetColumn(k)
		if col != nil {
			(*s.New)[k] = v
			switch col.TypeColumn {
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
* Insert
* @param data []et.Json
* @return *Command
**/
func (s *Model) Insert(data et.Json) *Command {
	return NewCommand(s, []et.Json{data}, Insert)
}

/**
* Update
* @param data []et.Json
* @return *Command
**/
func (s *Model) Update(data et.Json) *Command {
	return NewCommand(s, []et.Json{data}, Update)
}

/**
* Delete
* @return *Command
**/
func (s *Model) Delete() *Command {
	return NewCommand(s, []et.Json{}, Delete)
}

/**
* Bulk
* @param data []et.Json
* @return *Command
**/
func (s *Model) Bulk(data []et.Json) *Command {
	return NewCommand(s, data, Bulk)
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
		if len(s.Data) == 0 {
			return et.Items{}, mistake.New("Data not found")
		}

		for _, data := range s.Data {
			result, err := s.inserted(data)
			if err != nil {
				return et.Items{}, err
			}
			s.Result.Result = append(s.Result.Result, result.Result)
			s.Result.Count++
			s.Result.Ok = true
		}
	case Update:
		if len(s.Data) == 0 {
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
			result, err := s.updated(old, s.Data[0])
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

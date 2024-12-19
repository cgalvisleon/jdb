package jdb

import (
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
)

type TypeCommand int

const (
	Insert TypeCommand = iota
	Update
	Delete
)

type Command struct {
	Db      *DB
	Model   *Model
	Data    []et.Json
	Command TypeCommand
	Columns et.Json
	Atribs  et.Json
	New     *et.Json
	Key     string
	Wheres  []*LinqWhere
	Returns []*LinqSelect
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
		Command: command,
		Columns: et.Json{},
		Atribs:  et.Json{},
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
		if len(list) == 0 {
			return nil
		}

		from := &LinqFrom{
			Model: *s.Model,
			As:    s.Model.Table,
		}

		if len(list) == 1 {
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

func (s *Command) consolidate(data et.Json) et.Json {
	if s.Model.Integrity {
		for k, v := range data {
			if col := s.Model.GetColumn(k); col != nil {
				(*s.New)[k] = v
				switch col.TypeColumn {
				case TpAtribute:
					s.Atribs[k] = v
				case TpColumn:
					s.Columns[k] = v
				}
			}
		}
	} else {
		for k, v := range data {
			(*s.New)[k] = v
			if col := s.Model.GetColumn(k); col != nil {
				switch col.TypeColumn {
				case TpAtribute:
					s.Atribs[k] = v
				case TpColumn:
					s.Columns[k] = v
				}
			} else {
				s.Atribs[k] = v
			}
		}
	}

	return (*s.New)
}

func (s *Command) command() (et.Item, error) {
	result, err := s.Db.Command(s)
	if s.Show {
		logs.Debug(s.Describe().ToString())
	}
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}

func (s *Command) inserted(data et.Json) (et.Item, error) {
	s.consolidate(data)
	for _, trigger := range s.Model.BeforeInsert {
		err := trigger(nil, s.New, data)
		if err != nil {
			return et.Item{}, err
		}
	}

	result, err := s.command()
	if err != nil {
		return et.Item{}, err
	}

	for _, trigger := range s.Model.AfterInsert {
		err := trigger(nil, &result.Result, data)
		if err != nil {
			return et.Item{}, err
		}
	}

	s.Model.GetDetails(&result.Result)

	return result, nil
}

func (s *Command) updated(old, data et.Json) (et.Item, error) {
	s.consolidate(data)
	for _, trigger := range s.Model.BeforeUpdate {
		err := trigger(old, s.New, data)
		if err != nil {
			return et.Item{}, err
		}
	}

	result, err := s.command()
	if err != nil {
		return et.Item{}, err
	}

	if result.Ok {
		s.New = &result.Result
	}

	for _, trigger := range s.Model.AfterUpdate {
		err := trigger(old, &result.Result, data)
		if err != nil {
			return et.Item{}, err
		}
	}

	s.Model.GetDetails(&result.Result)

	return result, nil
}

func (s *Command) delete(old et.Json) (et.Item, error) {
	for _, trigger := range s.Model.BeforeDelete {
		err := trigger(old, nil, nil)
		if err != nil {
			return et.Item{}, err
		}
	}

	result, err := s.command()
	if err != nil {
		return et.Item{}, err
	}

	if result.Ok {
		s.New = &result.Result
	}

	for _, trigger := range s.Model.AfterDelete {
		err := trigger(old, nil, nil)
		if err != nil {
			return et.Item{}, err
		}
	}

	s.Model.GetDetails(&result.Result)

	return result, nil
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
	return NewCommand(s, data, Insert)
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

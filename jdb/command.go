package jdb

import (
	"database/sql"
	"slices"
	"strings"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

type TypeCommand int

const (
	Insert TypeCommand = iota
	Update
	Delete
	Bulk
	Sync
)

type Command struct {
	*QlWhere
	Tx      *sql.Tx
	Command TypeCommand
	Db      *DB
	From    *Model
	Data    []et.Json
	Values  []map[string]*Field
	Returns []*Field
	Sql     string
	Result  et.Items
	isUndo  bool
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
* ExecTx
* @param tx *sql.Tx
* @return et.Items, error
**/
func (s *Command) ExecTx(tx *sql.Tx) (et.Items, error) {
	switch s.Command {
	case Insert:
		if len(s.Data) == 0 {
			return et.Items{}, mistake.New(MSG_NOT_DATA)
		}

		err := s.inserted(tx)
		if err != nil {
			return et.Items{}, err
		}
	case Update:
		if len(s.Data) == 0 {
			return et.Items{}, mistake.New(MSG_NOT_DATA)
		}

		err := s.updated(tx)
		if err != nil {
			return et.Items{}, err
		}
	case Delete:
		err := s.deleted(tx)
		if err != nil {
			return et.Items{}, err
		}
	case Bulk:
		if len(s.Data) == 0 {
			return et.Items{}, mistake.New(MSG_NOT_DATA)
		}

		err := s.bulk(tx)
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
* @param tx *sql.Tx
* @return et.Item, error
**/
func (s *Command) OneTx(tx *sql.Tx) (et.Item, error) {
	result, err := s.ExecTx(tx)
	if err != nil {
		return et.Item{}, err
	}

	return result.First(), nil
}

/**
* Rollback
* @param tx *sql.Tx
* @return error
**/
func Rollback(tx *sql.Tx, err error) error {
	if tx == nil {
		return err
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		console.Error(mistake.Newf(MSG_ROLLBACK_ERROR, rollbackErr))
	}

	return err
}

/**
* Commit
* @param tx *sql.Tx
* @return error
**/
func Commit(tx *sql.Tx) error {
	if tx == nil {
		return nil
	}

	err := tx.Commit()
	if err != nil {
		console.Error(mistake.Newf(MSG_COMMIT_ERROR, err))
	}

	return err
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

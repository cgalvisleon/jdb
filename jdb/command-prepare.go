package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/timezone"
)

/**
* beforeInsertDefault
* @param tx *Tx, data et.Json
* @return error
**/
func (s *Command) beforeInsertDefault(tx *Tx, data et.Json) error {
	if s.From == nil {
		return mistake.New(MSG_MODEL_REQUIRED)
	}

	model := s.From

	if model.UseCore && model.IndexField != nil {
		data[model.IndexField.Name] = reg.GenIndex()
	}

	if model.UseCore && model.SystemKeyField != nil {
		data[model.SystemKeyField.Name] = model.GenId()
	}

	now := timezone.Now()
	if model.CreatedAtField != nil {
		data[model.CreatedAtField.Name] = now
	}

	if model.UpdatedAtField != nil {
		data[model.UpdatedAtField.Name] = now
	}

	for _, col := range model.Columns {
		if col.CalcFunction != nil {
			err := col.CalcFunction(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

/**
* beforeUpdateDefault
* @param tx *Tx, data et.Json
* @return error
**/
func (s *Command) beforeUpdateDefault(tx *Tx, data et.Json) error {
	if s.From == nil {
		return mistake.New(MSG_MODEL_REQUIRED)
	}

	model := s.From

	if model.CreatedAtField != nil {
		delete(data, model.CreatedAtField.Name)
	}

	if model.UpdatedAtField != nil {
		data[model.UpdatedAtField.Name] = timezone.Now()
	}

	for _, col := range model.Columns {
		if col.CalcFunction != nil {
			err := col.CalcFunction(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

/**
* prepare
* @return error
**/
func (s *Command) prepare() error {
	from := s.From
	s.Values = make([]map[string]*Field, 0)
	for i, data := range s.Data {
		value := make(map[string]*Field, 0)
		switch s.Command {
		case Insert:
			for _, fn := range s.beforeInsert {
				err := fn(s.tx, data)
				if err != nil {
					return err
				}

				s.Data[i] = data
			}
		case Update:
			for _, fn := range s.beforeUpdate {
				err := fn(s.tx, data)
				if err != nil {
					return err
				}

				s.Data[i] = data
			}
		case Delete:
			for _, fn := range s.beforeDelete {
				err := fn()
				if err != nil {
					return err
				}
			}
		}

		for k, v := range data {
			field := from.getField(k, true)
			if field == nil {
				continue
			}

			if field.Column == from.SourceField || field.Column == from.FullTextField {
				continue
			}

			field.setValue(v)

			if field.Column.TypeColumn != TpRelatedTo {
				value[field.Name] = field
			}
		}

		s.Values = append(s.Values, value)
	}

	if s.IsDebug {
		console.Debug(s.Describe().ToString())
	}

	return nil
}

package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/utility"
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

	if model.IndexField != nil && data.Int(model.IndexField.Name) == 0 {
		data[model.IndexField.Name] = reg.GenIndex()
	}

	if model.SystemKeyField != nil && data.Str(model.SystemKeyField.Name) == "" {
		data[model.SystemKeyField.Name] = model.GenId()
	}

	now := utility.Now()
	if model.CreatedAtField != nil && data.Str(model.CreatedAtField.Name) == "" {
		data[model.CreatedAtField.Name] = now
	}

	if model.UpdatedAtField != nil && data.Str(model.UpdatedAtField.Name) == "" {
		data[model.UpdatedAtField.Name] = now
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

	now := utility.Now()
	model := s.From
	if model.CreatedAtField != nil {
		delete(data, model.CreatedAtField.Name)
	}

	if model.UpdatedAtField != nil && data.Str(model.UpdatedAtField.Name) == "" {
		data[model.UpdatedAtField.Name] = now
	}

	return nil
}

/**
* prepare
* @return error
**/
func (s *Command) prepare() error {
	model := s.From
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
		}

		for k, v := range data {
			field := model.getField(k, true)
			if field == nil {
				continue
			}

			if field.Column == model.SourceField || field.Column == model.FullTextField {
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

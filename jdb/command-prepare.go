package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/timezone"
)

func (s *Command) beforeInsertDefault(data et.Json) (et.Json, error) {
	if s.From == nil {
		return data, mistake.New(MSG_MODEL_REQUIRED)
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
			for name, fn := range col.CalcFunction {
				val, err := fn(data)
				if err != nil {
					return data, err
				}

				data[name] = val
			}
		}
	}

	return data, nil
}

func (s *Command) beforeUpdateDefault(data et.Json) (et.Json, error) {
	if s.From == nil {
		return data, mistake.New(MSG_MODEL_REQUIRED)
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
			for name, fn := range col.CalcFunction {
				val, err := fn(data)
				if err != nil {
					return data, err
				}

				data[name] = val
			}
		}
	}

	return data, nil
}

func (s *Command) prepare() error {
	from := s.From
	s.Values = make([]map[string]*Field, 0)
	s.RelationsTo = make([]map[string]*Field, 0)
	for i, data := range s.Data {
		value := make(map[string]*Field, 0)
		relationsTo := make(map[string]*Field, 0)

		switch s.Command {
		case Insert:
			for _, fn := range s.beforeInsert {
				val, err := fn(data)
				if err != nil {
					return err
				}

				s.Data[i] = val
			}
		case Update:
			for _, fn := range s.beforeUpdate {
				val, err := fn(data)
				if err != nil {
					return err
				}

				s.Data[i] = val
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

			if field.Column.TypeColumn == TpRelatedTo {
				relationsTo[field.Name] = field
			} else {
				value[field.Name] = field
			}
		}

		s.Values = append(s.Values, value)
		s.RelationsTo = append(s.RelationsTo, relationsTo)
	}

	if s.IsDebug {
		console.Debug(s.Describe().ToString())
	}

	return nil
}

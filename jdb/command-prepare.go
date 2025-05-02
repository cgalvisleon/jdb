package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/timezone"
)

func (s *Command) beforeInsertDefault(data et.Json) et.Json {
	if s.From == nil {
		return data
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
				data[name] = fn(data)
			}
		}
	}

	return data
}

func (s *Command) beforeUpdateDefault(data et.Json) et.Json {
	if s.From == nil {
		return data
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
				data[name] = fn(data)
			}
		}
	}

	return data
}

func (s *Command) prepare() {
	from := s.From
	for i, data := range s.Data {
		value := make(map[string]*Field, 0)

		switch s.Command {
		case Insert:
			for _, fn := range s.beforeInsert {
				s.Data[i] = fn(data)
			}
		case Update:
			for _, fn := range s.beforeUpdate {
				s.Data[i] = fn(data)
			}
		}

		for k, v := range data {
			field := from.getField(k, true)
			if field == nil {
				continue
			}
			field.setValue(v)
			value[field.Name] = field
		}

		s.Values = append(s.Values, value)
	}
}

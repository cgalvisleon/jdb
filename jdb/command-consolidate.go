package jdb

import (
	"github.com/cgalvisleon/et/utility"
)

func setValue(value *Value, col *Column, v interface{}) *Value {
	switch col.TypeColumn {
	case TpAtribute:
		if value.Atribs[col.Name] == nil {
			value.Atribs[col.Name] = v
			value.Data[col.Name] = v
		}
	case TpColumn:
		if value.Columns[col.Name] == nil {
			value.Columns[col.Name] = v
			value.Data[col.Name] = v
		}
	}

	return value
}

func (s *Command) beforeInsert(value *Value) *Value {
	now := utility.Now()
	from := s.From
	if from.CreatedAtField != nil {
		setValue(value, from.CreatedAtField, now)
	}
	if from.UpdatedAtField != nil {
		setValue(value, from.UpdatedAtField, now)
	}
	if from.IndexField != nil {
		index := s.Db.GetSerie(from.Table)
		setValue(value, from.IndexField, index)
	}

	return value
}

func (s *Command) beforeUpdate(value *Value) *Value {
	now := utility.Now()
	from := s.From
	if from.UpdatedAtField != nil {
		setValue(value, from.UpdatedAtField, now)
	}

	return value
}

func (s *Command) consolidate() []*Value {
	from := s.From
	for _, data := range s.Origin {
		for k, v := range data {
			value := NewValue()
			field := from.GetField(k, true)
			if field != nil {
				setValue(value, field.Column, v)
			} else if from.SourceField != nil && !from.Integrity {
				value.Atribs[k] = v
				value.Data[k] = v
			}
		}
	}

	return s.Values
}

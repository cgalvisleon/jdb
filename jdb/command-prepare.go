package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/console"
)

/**
* prepare
* @return error
**/
func (s *Command) prepare() error {
	model := s.getModel()
	if model == nil {
		return fmt.Errorf(MSG_MODEL_REQUIRED)
	}

	s.Values = make([]map[string]*Field, 0)
	for i, data := range s.Data {
		value := make(map[string]*Field, 0)
		switch s.Command {
		case Insert:
			for _, fn := range model.beforeInsert {
				err := fn(s.tx, data)
				if err != nil {
					return err
				}
			}
			for _, fn := range s.beforeInsert {
				err := fn(s.tx, data)
				if err != nil {
					return err
				}

				s.Data[i] = data
			}
		case Update:
			for _, fn := range model.beforeUpdate {
				err := fn(s.tx, data)
				if err != nil {
					return err
				}
			}
			for _, fn := range s.beforeUpdate {
				err := fn(s.tx, data)
				if err != nil {
					return err
				}

				s.Data[i] = data
			}
		case Delete:
			for _, fn := range model.beforeDelete {
				err := fn(s.tx, data)
				if err != nil {
					return err
				}
			}
			for _, fn := range s.beforeDelete {
				err := fn(s.tx, data)
				if err != nil {
					return err
				}
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

package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/utility"
)

func (s *Command) beforeInsert(data et.Json) et.Json {
	now := utility.Now()
	from := s.From
	if from.CreatedAtField != nil && data[from.CreatedAtField.Name] == nil {
		data.Set(from.CreatedAtField.Name, now)
	}
	if from.UpdatedAtField != nil && data[from.UpdatedAtField.Name] == nil {
		data.Set(from.UpdatedAtField.Name, now)
	}
	if from.IndexField != nil && data[from.IndexField.Name] == nil {
		index := s.Db.GetSerie(from.Table)
		data.Set(from.IndexField.Name, index)
	}

	return data
}

func (s *Command) inserted(data et.Json) (et.Item, error) {
	data = s.beforeInsert(data)
	s.consolidate(data)

	result, err := s.Db.Command(s)
	if err != nil {
		return et.Item{}, err
	}

	if !result.Ok {
		return et.Item{}, mistake.New(MSG_NOT_INSERT_DATA)
	}

	model := s.From.Model
	before := result.Json(0, "before")
	after := result.Json(0, "after")

	for _, event := range s.From.EventsInsert {
		err := event(model, before, &after, data)
		if err != nil {
			return et.Item{}, err
		}
	}

	return et.Item{
		Ok: true,
		Result: et.Json{
			"before": before,
			"after":  after,
		},
	}, nil
}

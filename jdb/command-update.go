package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

func (s *Command) beforeUpdate(data et.Json) et.Json {
	now := utility.Now()
	from := s.From
	if from.UpdatedAtField != nil && data[from.UpdatedAtField.Name] == nil {
		data.Set(from.UpdatedAtField.Name, now)
	}

	return data
}

func (s *Command) updated(data et.Json) (et.Items, error) {
	data = s.beforeInsert(data)
	s.consolidate(data)

	results, err := s.Db.Command(s)
	if err != nil {
		return et.Items{}, err
	}

	model := s.From.Model
	for _, result := range results.Result {
		before := result.Json("before")
		after := result.Json("after")

		for _, event := range s.From.EventsUpdate {
			err := event(model, before, &after, data)
			if err != nil {
				return et.Items{}, err
			}
		}
	}

	return results, nil
}

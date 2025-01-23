package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

type UndoRecord struct {
	Key   string `json:"key"`
	Index int64  `json:"index"`
}

func (s *Command) undo() error {
	if s.Undo == nil {
		return mistake.New(MSG_UNDO_NOT_DEFINED)
	}

	from := s.From
	history := from.History
	if history == nil {
		return mistake.New(MSG_HISTORY_NOT_DEFINED)
	}

	historyData, err := history.
		Where(from.KeyField.Name).Eq(s.Undo.Key).
		And(HISTORY_INDEX).Eq(s.Undo.Index).
		One()
	if err != nil {
		return err
	}

	if !historyData.Ok {
		return mistake.New(MSG_HISTORY_NOT_FOUND)
	}

	s.Origin = []et.Json{historyData.Result}
	err = s.update()
	if err != nil {
		return err
	}

	if !s.Result.Ok {
		return mistake.New(MSG_NOT_UPDATE_DATA)
	}

	go history.Delete().
		Where(from.KeyField.Name).Eq(s.Undo.Key).
		And(HISTORY_INDEX).More(s.Undo.Index).
		Exec()

	for _, result := range s.Result.Result {
		before := result.Json("before")
		after := result.Json("after")

		for _, event := range s.From.EventsUpdate {
			err := event(from.Model, before, after)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

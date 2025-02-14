package jdb

import (
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
	if history.With == nil {
		return mistake.New(MSG_HISTORY_NOT_DEFINED)
	}

	old, err := history.With.
		Where(history.Fk.Name).Eq(s.Undo.Key).
		And(HISTORY_INDEX).Eq(s.Undo.Index).
		One()
	if err != nil {
		return err
	}

	if !old.Ok {
		return mistake.New(MSG_HISTORY_NOT_FOUND)
	}

	model := from.Model
	_, err = model.Update(old.Result).
		Where(history.Fk.Name).Eq(s.Undo.Key).
		History(false).
		Exec()
	if err != nil {
		return err
	}

	go history.With.Delete().
		Where(history.Fk.Name).Eq(s.Undo.Key).
		And(HISTORY_INDEX).More(s.Undo.Index).
		Exec()

	return nil
}

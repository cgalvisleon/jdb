package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) undo() error {
	if s.Undo.IsEmpty() {
		return mistake.New(MSG_UNDO_NOT_DEFINED)
	}

	key := s.Undo.ValStr("", "key")
	if key == "" {
		return mistake.New(MSG_UNDO_KEY_NOT_DEFINED)
	}

	history := s.From.History.With
	if history == nil {
		return mistake.New(MSG_HISTORY_NOT_DEFINED)
	}

	from := s.From
	if len(from.PrimaryKeys) == 0 {
		return mistake.New(MSG_PRIMARYKEY_NOT_FOUND)
	}

	pkn := from.PrimaryKeys[0].Name

	var err error
	var old et.Item
	index := s.Undo.ValInt64(-1, "index")
	if index == -1 {
		old, err = history.
			Where(pkn).Eq(key).
			OrderByDesc(HISTORY_INDEX).
			One()
		if err != nil {
			return err
		}

		index = old.Int64(HISTORY_INDEX)
	} else {
		old, err = history.
			Where(pkn).Eq(key).
			And(HISTORY_INDEX).Eq(index).
			One()
		if err != nil {
			return err
		}
	}

	if !old.Ok {
		return mistake.New(MSG_HISTORY_NOT_FOUND)
	}

	delete(old.Result, HISTORY_INDEX)
	model := s.From
	_, err = model.Update(old.Result).
		Where(pkn).Eq(key).
		History(false).
		Exec()
	if err != nil {
		return err
	}

	go history.Delete().
		Where(pkn).Eq(key).
		And(HISTORY_INDEX).More(index).
		Exec()

	return nil
}

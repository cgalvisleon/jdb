package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) current(where et.Json) error {
	if len(s.Data) != 1 {
		return mistake.New(MSG_MANY_INSERT_DATA)
	}

	model := s.From
	columns := model.getColumnsByType(TpColumn)
	mainWhere := s.getWheres()
	ql := From(model)
	ql.SetWheres(where)
	ql.SetWheres(mainWhere)
	ql.SetSelects(columns)
	current, err := ql.
		SetDebug(s.IsDebug).
		AllTx(s.tx)
	if err != nil {
		return err
	}

	s.Current = current
	mapCurrent, err := model.getMapByPk(current.Result)
	if err != nil {
		return err
	}

	s.CurrentMap = mapCurrent

	return nil
}

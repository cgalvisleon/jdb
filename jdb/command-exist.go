package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

func (s *Command) exists(model *Model, where et.Json) (bool, error) {
	if model == nil {
		return false, mistake.New(MSG_MODEL_REQUIRED)
	}

	ql := From(model)
	ql.setWheres(where)
	exist, err := ql.
		setDebug(s.IsDebug).
		ItExistsTx(s.tx)
	if err != nil {
		return false, err
	}

	return exist, nil
}

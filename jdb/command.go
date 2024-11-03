package jdb

import (
	"github.com/cgalvisleon/et/et"
)

type Command struct{}

func (s *Model) Insert(data et.Json) *Command {
	return &Command{}
}

func (s *Model) Update(data et.Json) *Command {
	return &Command{}
}

func (s *Model) Delete(_where LinqWhere) *Command {
	return &Command{}
}

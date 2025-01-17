package jdb

import (
	"github.com/cgalvisleon/et/et"
)

type TypeTrigger int

const (
	BeforeInsert TypeTrigger = iota
	AfterInsert
	BeforeUpdate
	AfterUpdate
	BeforeDelete
	AfterDelete
)

func (s TypeTrigger) Name() string {
	switch s {
	case BeforeInsert:
		return "BeforeInsert"
	case AfterInsert:
		return "AfterInsert"
	case BeforeUpdate:
		return "BeforeUpdate"
	case AfterUpdate:
		return "AfterUpdate"
	case BeforeDelete:
		return "BeforeDelete"
	case AfterDelete:
		return "AfterDelete"
	}
	return ""
}

type Trigger func(old et.Json, new *et.Json, data et.Json) error

var Triggers = map[string]Trigger{}

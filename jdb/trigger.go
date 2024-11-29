package jdb

import "github.com/cgalvisleon/et/et"

type TypeTrigger int

const (
	BeforeInsert TypeTrigger = iota
	AfterInsert
	BeforeUpdate
	AfterUpdate
	BeforeDelete
	AfterDelete
)

type Trigger func(old et.Json, new *et.Json, data et.Json) error

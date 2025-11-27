package jdb

import "github.com/cgalvisleon/et/et"

type Detail struct {
	From    *Model   `json:"from"`
	Fks     et.Json  `json:"fks"`
	Selects []string `json:"selects"`
}

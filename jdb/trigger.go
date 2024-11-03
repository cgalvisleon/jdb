package jdb

import "github.com/cgalvisleon/et/et"

type Trigger func(model *Model, old, new, data *et.Json) error

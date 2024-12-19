package jdb

import "github.com/cgalvisleon/et/et"

type HandlerListener func(res et.Json)

var listenerChannels map[string]bool = map[string]bool{}

func (s *Model) OnListener(channels []string, listener HandlerListener) {

}

package jdb

import "github.com/cgalvisleon/et/et"

type Column struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Default interface{} `json:"default"`
	Hidden  bool        `json:"hidden"`
}

func (s *Column) ToJson() et.Json {
	return et.Json{
		"name":    s.Name,
		"type":    s.Type,
		"default": s.Default,
		"hidden":  s.Hidden,
	}
}

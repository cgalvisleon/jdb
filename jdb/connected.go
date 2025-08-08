package jdb

import "github.com/cgalvisleon/et/et"

type Connected interface {
	Chain() (string, error)
	ToJson() et.Json
	Validate() error
}

type ConnectParams struct {
	Id       string    `json:"id"`
	Driver   string    `json:"driver"`
	Name     string    `json:"name"`
	UserCore bool      `json:"user_core"`
	NodeId   int       `json:"node_id"`
	Debug    bool      `json:"debug"`
	Params   Connected `json:"params"`
}

/**
* Json
* @return et.Json
**/
func (s *ConnectParams) ToJson() et.Json {
	return et.Json{
		"id":        s.Id,
		"driver":    s.Driver,
		"name":      s.Name,
		"user_core": s.UserCore,
		"node_id":   s.NodeId,
		"params":    s.Params.ToJson(),
	}
}

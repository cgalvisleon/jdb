package jdb

import (
	"encoding/json"

	"github.com/cgalvisleon/et/et"
)

const (
	TpRows   = "rows"
	TpObject = "object"
)

var (
	TpQuerys = map[string]bool{
		TpObject: true,
		TpRows:   true,
	}
)

type Ql struct {
	*where
	Database  string                  `json:"database"`
	Type      string                  `json:"type"`
	Froms     et.Json                 `json:"from"`
	Selects   et.Json                 `json:"selects"`
	Atribs    et.Json                 `json:"atribs"`
	Hidden    []string                `json:"hidden"`
	Rollups   et.Json                 `json:"rollups"`
	Relations et.Json                 `json:"relations"`
	Calls     et.Json                 `json:"calls"`
	Joins     []et.Json               `json:"joins"`
	GroupBy   []string                `json:"group_by"`
	Havings   []et.Json               `json:"having"`
	OrderBy   et.Json                 `json:"order_by"`
	Limits    et.Json                 `json:"limit"`
	Exists    bool                    `json:"exists"`
	Count     bool                    `json:"count"`
	SQL       string                  `json:"sql"`
	calls     map[string]*DataContext `json:"-"`
	db        *Database               `json:"-"`
	tx        *Tx                     `json:"-"`
	isDebug   bool                    `json:"-"`
}

/**
* NewQl
* @return *Ql
**/
func newQl(db *Database) *Ql {
	return &Ql{
		where:     newWhere(),
		Database:  db.Name,
		Type:      TpRows,
		Froms:     et.Json{},
		Selects:   et.Json{},
		Atribs:    et.Json{},
		Hidden:    []string{},
		Rollups:   et.Json{},
		Relations: et.Json{},
		Calls:     et.Json{},
		Joins:     make([]et.Json, 0),
		GroupBy:   []string{},
		Havings:   make([]et.Json, 0),
		OrderBy:   et.Json{},
		Limits:    et.Json{},
		db:        db,
		calls:     make(map[string]*DataContext),
	}
}

/**
* ToJson
* @return et.Json
**/
func (s *Ql) ToJson() et.Json {
	bt, err := json.Marshal(s)
	if err != nil {
		return et.Json{}
	}

	var result et.Json
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* Debug
* @return *Ql
**/
func (s *Ql) Debug() *Ql {
	s.isDebug = true
	return s
}

/**
* addFrom
* @param name, as string
* @return *Ql
**/
func (s *Ql) addFrom(name, as string) *Ql {
	s.Froms[name] = as
	return s
}

/**
* setTx
* @param tx *Tx
* @return *Ql
**/
func (s *Ql) setTx(tx *Tx) *Ql {
	s.tx = tx
	return s
}

/**
* Select
* @param fields interface{}
* @return *Ql
**/
func (s *Ql) Select(fields interface{}) *Ql {
	switch v := fields.(type) {
	case string:
		s.Selects[v] = v
	case []string:
		for _, v := range v {
			s.Selects[v] = v
		}
	case et.Json:
		for k, v := range v {
			s.Selects[k] = v
		}
	}

	return s
}

/**
* Object
* @param fields interface{}
* @return *Ql
**/
func (s *Ql) Object(fields interface{}) *Ql {
	switch v := fields.(type) {
	case string:
		s.Selects[v] = v
	case []string:
		for _, v := range v {
			s.Selects[v] = v
		}
	case et.Json:
		for k, v := range v {
			s.Selects[k] = v
		}
	}

	s.Type = TpObject
	return s
}

/**
* Join
* @param to, as string, on []Condition
* @return *Ql
*
 */
func (s *Ql) Join(to, as string, on ...Condition) *Ql {
	n := len(s.Froms)
	if n == 0 {
		return s
	}

	model, err := s.db.getModel(to)
	if err != nil {
		return s
	}

	ons := make([]et.Json, 0)
	for _, v := range on {
		ons = append(ons, v.ToJson())
	}

	n = len(s.Joins) + 1
	s.Joins = append(s.Joins, et.Json{
		"from": et.Json{
			model.Table: as,
		},
		"on": ons,
	})

	return s
}

/**
* Group
* @param fields ...string
* @return *Ql
**/
func (s *Ql) Group(fields ...string) *Ql {
	s.GroupBy = append(s.GroupBy, fields...)
	return s
}

/**
* Having
* @param having et.Json
* @return *Ql
**/
func (s *Ql) Having(having ...Condition) *Ql {
	havings := make([]et.Json, 0)
	for _, v := range having {
		havings = append(havings, v.ToJson())
	}

	s.Havings = havings
	return s
}

/**
* Order
* @param asc bool, fields ...string
* @return *Ql
**/
func (s *Ql) Order(asc bool, fields ...string) *Ql {
	if asc {
		s.OrderBy["asc"] = fields
	} else {
		s.OrderBy["desc"] = fields
	}

	return s
}

/**
* Page
* @param page, rows int
* @return *Ql
**/
func (s *Ql) Limit(page, rows int) (items et.Items, err error) {
	s.Limits["page"] = page
	s.Limits["rows"] = rows
	return s.queryTx(nil)
}

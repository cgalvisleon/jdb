package jdb

import (
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/et"
)

type Ql struct {
	*where
	Database    string                 `json:"database"`
	SourceField string                 `json:"source_field"`
	Froms       et.Json                `json:"from"`
	Selects     et.Json                `json:"selects"`
	Atribs      et.Json                `json:"atribs"`
	Hiddens     []string               `json:"hidden"`
	Calcs       map[string]DataContext `json:"-"`
	Vms         map[string]string      `json:"vms"`
	Rollups     et.Json                `json:"rollups"`
	Relations   et.Json                `json:"relations"`
	Joins       []et.Json              `json:"joins"`
	GroupBy     []string               `json:"group_by"`
	Havings     []et.Json              `json:"having"`
	OrdersBy    et.Json                `json:"order_by"`
	Limits      et.Json                `json:"limit"`
	Exists      bool                   `json:"exists"`
	Count       bool                   `json:"count"`
	SQL         string                 `json:"sql"`
	db          *DB                    `json:"-"`
	tx          *Tx                    `json:"-"`
	IsDebug     bool                   `json:"-"`
	useJoin     bool                   `json:"-"`
}

/**
* NewQl
* @return *Ql
**/
func newQl(model *Model, as string) *Ql {
	result := &Ql{
		where:     newWhere(model, as),
		Database:  model.Database,
		Froms:     et.Json{},
		Selects:   et.Json{},
		Atribs:    et.Json{},
		Hiddens:   []string{},
		Rollups:   et.Json{},
		Relations: et.Json{},
		Calcs:     make(map[string]DataContext),
		Joins:     make([]et.Json, 0),
		GroupBy:   []string{},
		Havings:   make([]et.Json, 0),
		OrdersBy:  et.Json{},
		Limits:    et.Json{},
		db:        model.db,
	}

	if model != nil {
		result.addFrom(model, as)
	}

	return result
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
	s.setDebug(true)
	return s
}

/**
* addFrom
* @param model *Model, as string
* @return *Ql
**/
func (s *Ql) addFrom(model *Model, as string) *Ql {
	s.From = model
	s.As = as
	s.Froms[model.Table] = as
	if s.SourceField == "" {
		s.SourceField = model.SourceField
	}

	for _, v := range model.Hidden {
		v = fmt.Sprintf("%s.%s", as, v)
		s.Hiddens = append(s.Hiddens, v)
	}
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
* Page
* @param page int
* @return *Ql
**/
func (s *Ql) Page(page int) *Ql {
	s.Limits["page"] = page
	return s
}

/**
* Hidden
* @param hiddens ...string
* @return *Ql
**/
func (s *Ql) Hidden(hiddens ...string) *Ql {
	s.Hiddens = hiddens
	return s
}

/**
* QueryTx
* @param tx *Tx, query et.Json
* @return et.Items, error
**/
func (s *Ql) QueryTx(tx *Tx, query et.Json) (et.Items, error) {
	s.setQuery(query)
	return s.queryTx(tx)
}

/**
* Query
* @param query et.Json
* @return et.Items, error
**/
func (s *Ql) Query(query et.Json) (et.Items, error) {
	return s.QueryTx(nil, query)
}

package jdb

import (
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
)

var MaxRows int

func init() {
	MaxRows = envar.GetInt("MAX_ROWS", 1000)
}

type Ql struct {
	*where
	Database  string                 `json:"database"`
	Selects   et.Json                `json:"selects"`
	Atribs    et.Json                `json:"atribs"`
	Hiddens   []string               `json:"hidden"`
	Details   map[string]*Detail     `json:"details"`
	Rollups   map[string]*Detail     `json:"rollups"`
	Relations map[string]*Detail     `json:"relations"`
	Calcs     map[string]DataContext `json:"-"`
	Joins     []et.Json              `json:"joins"`
	Wheres    []et.Json              `json:"where"`
	GroupBy   []string               `json:"group_by"`
	Havings   []et.Json              `json:"having"`
	OrdersBy  et.Json                `json:"order_by"`
	Limits    et.Json                `json:"limit"`
	Exists    bool                   `json:"exists"`
	Count     bool                   `json:"count"`
	SQL       string                 `json:"sql"`
	MaxRows   int                    `json:"max_rows"`
	IsDebug   bool                   `json:"-"`
	db        *DB                    `json:"-"`
	tx        *Tx                    `json:"-"`
	inJoin    bool                   `json:"-"`
}

/**
* NewQl
* @return *Ql
**/
func newQl(model *Model, as string) *Ql {
	result := &Ql{
		where:     newWhere(model, as),
		Database:  model.Database,
		Selects:   et.Json{},
		Atribs:    et.Json{},
		Hiddens:   []string{},
		Details:   make(map[string]*Detail),
		Rollups:   make(map[string]*Detail),
		Relations: make(map[string]*Detail),
		Calcs:     make(map[string]DataContext),
		Joins:     make([]et.Json, 0),
		Wheres:    make([]et.Json, 0),
		GroupBy:   []string{},
		Havings:   make([]et.Json, 0),
		OrdersBy:  et.Json{},
		MaxRows:   MaxRows,
		Limits:    et.Json{},
		db:        model.db,
	}

	if model != nil {
		result.addFroms(model, as)
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
func (s *Ql) addFroms(model *Model, as string) *Ql {
	s.Froms[as] = model

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

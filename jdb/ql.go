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
	UseAtribs   bool                   `json:"use_atribs"`
	SQL         string                 `json:"sql"`
	db          *DB                    `json:"-"`
	from        *Model                 `json:"-"`
	tx          *Tx                    `json:"-"`
	isDebug     bool                   `json:"-"`
	useJoin     bool                   `json:"-"`
}

/**
* NewQl
* @return *Ql
**/
func newQl(db *DB) *Ql {
	return &Ql{
		where:     newWhere(),
		Database:  db.Name,
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
		db:        db,
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
* @param model *Model, as string
* @return *Ql
**/
func (s *Ql) addFrom(model *Model, as string) *Ql {
	s.from = model
	s.Froms[model.Table] = as
	s.SourceField = model.SourceField
	s.UseAtribs = s.SourceField != ""
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
func (s *Ql) Select(fields ...string) *Ql {
	if s.from == nil {
		return s
	}

	for _, v := range fields {
		col, ok := s.from.GetColumn(v)
		tp := col.String("type")
		if ok && TypeColumn[tp] {
			s.Selects[v] = v
			continue
		}

		if s.UseAtribs || TypeAtrib[tp] {
			s.Atribs[v] = v
			continue
		}

		if tp == TypeCalc {
			s.Calcs[v] = s.from.Calcs[v]
		} else if tp == TypeVm {
			s.Vms[v] = s.from.Vms[v]
		} else if tp == TypeRollup {
			s.Rollups[v] = s.from.Rollups[v]
		} else if tp == TypeRelation {
			s.Relations[v] = s.from.Relations[v]
		} else if tp == TypeDetail {
			to, err := s.db.GetModel(v)
			if err != nil {
				continue
			}

			detail := s.from.Details[v]
			references := detail.Json("references")
			columns := references.ArrayJson("columns")
			as := string(rune(len(s.Joins) + 66))
			first := true
			for _, fk := range columns {
				for k, v := range fk {
					if first {
						s.Join(to, as, Eq(fmt.Sprintf("A.%s", k), fmt.Sprintf("%s.%s", as, v)))
						first = false
						continue
					}
					s.And(Eq(fmt.Sprintf("A.%s", k), fmt.Sprintf("%s.%s", as, v)))
				}
			}
		}
	}

	return s
}

/**
* Join
* @param to *Model, as string, on Condition
* @return *Ql
*
 */
func (s *Ql) Join(to *Model, as string, on Condition) *Ql {
	n := len(s.Froms)
	if n == 0 {
		return s
	}

	n = len(s.Joins) + 1
	s.Joins = append(s.Joins, et.Json{
		"from": et.Json{
			to.Table: as,
		},
		"on": []et.Json{
			on.ToJson(),
		},
	})
	s.useJoin = true

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
		s.OrdersBy["asc"] = fields
	} else {
		s.OrdersBy["desc"] = fields
	}

	return s
}

/**
* OrderBy
* @param fields ...string
* @return *Ql
**/
func (s *Ql) OrderBy(fields ...string) *Ql {
	s.OrdersBy["asc"] = fields
	return s
}

/**
* OrderDesc
* @param fields ...string
* @return *Ql
**/
func (s *Ql) OrderDesc(fields ...string) *Ql {
	s.OrdersBy["desc"] = fields
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

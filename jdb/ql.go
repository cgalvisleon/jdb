package jdb

import (
	"encoding/json"

	"github.com/cgalvisleon/et/et"
)

type Ql struct {
	Database  string                  `json:"-"`
	Froms     []string                `json:"froms"`
	Selects   et.Json                 `json:"selects"`
	Atributes et.Json                 `json:"atributes"`
	Rollups   et.Json                 `json:"rollups"`
	Relations et.Json                 `json:"relations"`
	Joins     []et.Json               `json:"joins"`
	Wheres    et.Json                 `json:"wheres"`
	GroupBy   et.Json                 `json:"group_by"`
	Having    et.Json                 `json:"having"`
	OrderBy   et.Json                 `json:"order_by"`
	Limit     et.Json                 `json:"limit"`
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
		Database:  db.Name,
		Froms:     []string{},
		Selects:   et.Json{},
		Atributes: et.Json{},
		Rollups:   et.Json{},
		Relations: et.Json{},
		Joins:     make([]et.Json, 0),
		Wheres:    et.Json{},
		GroupBy:   et.Json{},
		Having:    et.Json{},
		OrderBy:   et.Json{},
		Limit:     et.Json{},
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
* getAs
* @return string
**/
func getAs(ql *Ql) string {
	n := len(ql.Froms)
	as := string(rune(65 + n))
	return as
}

/**
* addFrom
* @param name string
* @return *Ql
**/
func (s *Ql) addFrom(name string) *Ql {
	s.Froms = append(s.Froms, name)
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

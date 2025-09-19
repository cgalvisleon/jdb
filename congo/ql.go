package jdb

import (
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

type Ql struct {
	Froms   []string  `json:"froms"`
	Selects []string  `json:"selects"`
	Joins   []et.Json `json:"joins"`
	Wheres  et.Json   `json:"wheres"`
	GroupBy et.Json   `json:"group_by"`
	Having  et.Json   `json:"having"`
	OrderBy et.Json   `json:"order_by"`
	Limit   et.Json   `json:"limit"`
	SQL     string    `json:"sql"`
	db      *Database `json:"-"`
	tx      *Tx       `json:"-"`
	isDebug bool      `json:"-"`
}

/**
* NewQl
* @return *Ql
**/
func newQl(db *Database) *Ql {
	return &Ql{
		Froms:   []string{},
		Selects: []string{},
		Joins:   make([]et.Json, 0),
		Wheres:  et.Json{},
		GroupBy: et.Json{},
		Having:  et.Json{},
		OrderBy: et.Json{},
		Limit:   et.Json{},
		db:      db,
	}
}

/**
* Query
* @param query et.Json
* @return (*Ql, error)
**/
func Query(query et.Json) (*Ql, error) {
	database := query.String("database")
	if !utility.ValidStr(database, 0, []string{}) {
		return nil, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	db, err := GetDatabase(database)
	if err != nil {
		return nil, err
	}

	result := newQl(db)
	from := query.Str("from")
	if len(from) == 0 {
		return nil, fmt.Errorf(MSG_FROM_REQUIRED)
	}
	result.addFrom(from)

	return result.setQuery(query), nil
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
func (s *Ql) getAs() string {
	n := len(s.Joins)
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

package postgres

import (
	"database/sql"
	"encoding/json"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) defineFlows() error {
	exist, err := s.existTable("core", "FLOWS")
	if err != nil {
		return console.Panic(err)
	}

	if exist {
		return nil
	}

	sql := parceSQL(`
  CREATE TABLE IF NOT EXISTS core.FLOWS(
		_ID VARCHAR(80) DEFAULT '',
		NAME VARCHAR(250) DEFAULT '',
		DESCRIPTION TEXT DEFAULT '',
		VALUE BYTEA,
		_IDT VARCHAR(80) DEFAULT '-1',
		INDEX BIGINT DEFAULT 0,
		PRIMARY KEY(_ID)
	);
	CREATE INDEX IF NOT EXISTS FLOWS_NAME_IDX ON core.FLOWS(NAME);
	CREATE INDEX IF NOT EXISTS FLOWS__IDT_IDX ON core.FLOWS(_IDT);
	CREATE INDEX IF NOT EXISTS FLOWS_INDEX_IDX ON core.FLOWS(INDEX);`)
	sql = strs.Append(sql, defineRecordTrigger("core.FLOWS"), "\n")
	sql = strs.Append(sql, defineSeriesTrigger("core.FLOWS"), "\n")

	err = s.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

/**
* SetFlow
* @params name string
* @params value []byte
* @return error
**/
func (s *Postgres) SetFlow(name string, value []byte) error {
	sql := parceSQL(`
	UPDATE core.FLOWS SET
	VALUE = $2
	WHERE NAME = $1
	RETURNING *;`)

	items, err := s.Query(sql, name, value)
	if err != nil {
		return err
	}

	if items.Ok {
		return nil
	}

	sql = parceSQL(`
	INSERT INTO core.FLOWS(_ID, NAME, VALUE, INDEX)
	VALUES ($1, $2, $3, $4);`)

	id := utility.RecordId("flow", "")
	index := s.GetSerie("core.FLOWS")
	err = s.Exec(sql, id, name, value, index)
	if err != nil {
		return err
	}

	return nil
}

/**
* GetFlow
* @params id string
* @return jdb.Flow, error
**/
func (s *Postgres) GetFlow(id string) (jdb.Flow, error) {
	query := parceSQL(`
	SELECT VALUE, INDEX
	FROM core.FLOWS
	WHERE _ID = $1
	LIMIT 1;`)

	var ok bool
	var value []byte
	var index int
	err := s.db.QueryRow(query, id).Scan(&value, &index)
	if err != nil {
		ok = err == sql.ErrNoRows
		if !ok {
			return jdb.Flow{}, err
		}
	}

	var result jdb.Flow
	err = json.Unmarshal(value, &result)
	if err != nil {
		return jdb.Flow{}, err
	}

	return result, nil
}

/**
* DeleteFlow
* @params id string
* @return error
**/
func (s *Postgres) DeleteFlow(id string) error {
	sql := parceSQL(`
	DELETE 
	FROM core.FLOWS
	WHERE _ID = $1;`)

	err := s.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}

/**
* FindFlows - Find key value
* @params search string
* @params page int
* @params rows int
* @return et.List, error
**/
func (s *Postgres) FindFlows(search string, page, rows int) (et.List, error) {
	sql := `
	SELECT COUNT(*) AS ALL
	FROM core.FLOWS A
	WHERE A.NAME ILIKE %$1%;`

	result, err := s.Query(sql, search)
	if err != nil {
		return et.List{}, err
	}

	all := result.Int(0, "all")

	sql = `
	SELECT A.VALUE, A.INDEX
	FROM core.FLOWS A
	WHERE A.NAME ILIKE %$1%
	OFFSET $2 LIMIT $3
	ORDER BY A.INDEX;`

	offset := (page - 1) * rows
	rws, err := s.db.Query(sql, search, offset, rows)
	if err != nil {
		return et.List{}, err
	}
	defer rws.Close()

	flows := et.Items{}
	for rws.Next() {
		var value []byte
		var index int
		err = rws.Scan(&value, &index)
		if err != nil {
			return et.List{}, err
		}

		var flow jdb.Flow
		err = json.Unmarshal(value, &flow)
		if err != nil {
			return et.List{}, err
		}

		flows.Add(flow.Describe())
	}

	return flows.ToList(all, page, rows), nil
}

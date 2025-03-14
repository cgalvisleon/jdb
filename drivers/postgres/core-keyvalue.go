package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* defineKeyValue - Define KeyValue table
* @return error
**/
func (s *Postgres) defineKeyValue() error {
	exist, err := s.existTable("core", "KEYVALUES")
	if err != nil {
		return console.Panic(err)
	}

	if exist {
		return nil
	}

	sql := parceSQL(`
  CREATE TABLE IF NOT EXISTS core.KEYVALUES(
		_ID VARCHAR(80) DEFAULT '',
		VALUE BYTEA,
		_IDT VARCHAR(80) DEFAULT '-1',
		INDEX BIGINT DEFAULT 0,
		PRIMARY KEY(_ID)
	);
	CREATE INDEX IF NOT EXISTS KEYVALUES__IDT_IDX ON core.KEYVALUES(_IDT);
	CREATE INDEX IF NOT EXISTS KEYVALUES_INDEX_IDX ON core.KEYVALUES(INDEX);`)
	sql = strs.Append(sql, defineRecordTrigger("core.KEYVALUES"), "\n")
	sql = strs.Append(sql, defineSeriesTrigger("core.KEYVALUES"), "\n")

	err = s.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

/**
* SetKey - Set key value
* @params key string
* @params value []byte
* @return error
**/
func (s *Postgres) SetKey(key string, value []byte) error {
	sql := parceSQL(`
	UPDATE core.KEYVALUES SET
	VALUE = $2
	WHERE _ID = $1
	RETURNING *;`)

	items, err := s.Query(sql, key, value)
	if err != nil {
		return err
	}

	if items.Ok {
		return nil
	}

	sql = parceSQL(`
	INSERT INTO core.KEYVALUES(_ID, VALUE, INDEX)
	VALUES ($1, $2, $3)
	ON CONFLICT (_ID) DO UPDATE SET
		VALUE = $2
	WHERE
		core.KEYVALUES.VALUE IS DISTINCT FROM EXCLUDED.VALUE;`)

	index := s.GetSerie("core.KEYVALUES")
	err = s.Exec(sql, key, value, index)
	if err != nil {
		return err
	}

	return nil
}

/**
* GetKey - Get key value
* @params key string
* @return KeyValue, error
**/
func (s *Postgres) GetKey(key string) (et.KeyValue, error) {
	sql := parceSQL(`
	SELECT VALUE, INDEX
	FROM core.KEYVALUES
	WHERE _ID = $1
	LIMIT 1;`)

	var ok bool
	var value []byte
	var index int
	sql = jdb.SQLParse(sql, key)
	ok, err := s.QueryRow(sql, &value, &index)
	if err != nil {
		return et.KeyValue{}, err
	}

	return et.KeyValue{
		Ok:    !ok,
		Value: value,
		Imdex: index,
	}, nil
}

/**
* DeleteKey - Delete key value
* @params key string
* @return error
**/
func (s *Postgres) DeleteKey(key string) error {
	sql := parceSQL(`
	DELETE 
	FROM core.KEYVALUES
	WHERE _ID = $1;`)

	err := s.Exec(sql, key)
	if err != nil {
		return err
	}

	return nil
}

/**
* FindKeys - Find key value
* @params search string
* @params page int
* @params rows int
* @return et.List, error
**/
func (s *Postgres) FindKeys(search string, page, rows int) (et.List, error) {
	sql := `
	SELECT COUNT(*) AS ALL
	FROM core.KEYVALUES A
	WHERE A.VALUE ILIKE %$1%;`

	result, err := s.Query(sql, search)
	if err != nil {
		return et.List{}, err
	}

	all := result.Int(0, "all")

	sql = `
	SELECT A.*
	FROM core.KEYVALUES A
	WHERE A.VALUE ILIKE %$1%
	OFFSET $2 LIMIT $3
	ORDER BY A.INDEX;`

	offset := (page - 1) * rows
	items, err := s.Query(sql, search, offset, rows)
	if err != nil {
		return et.List{}, err
	}

	return items.ToList(all, page, rows), nil
}

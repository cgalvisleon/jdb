package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
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

	sql := `
  CREATE TABLE IF NOT EXISTS core.KEYVALUES(
		_ID VARCHAR(80) DEFAULT '',
		VALUE TEXT,
		_IDT VARCHAR(80) DEFAULT '-1',
		INDEX BIGINT DEFAULT 0,
		PRIMARY KEY(_ID)
	);
	CREATE INDEX IF NOT EXISTS KEYVALUES__IDT_IDX ON core.KEYVALUES(_IDT);
	CREATE INDEX IF NOT EXISTS KEYVALUES_INDEX_IDX ON core.KEYVALUES(INDEX);
	`
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
* @params value string
* @return error
**/
func (s *Postgres) SetKey(key, value string) error {
	sql := `
	UPDATE core.KEYVALUES SET
	VALUE = $2
	WHERE _ID = $1
	RETURNING *;`

	item, err := s.One(sql, key, value)
	if err != nil {
		return err
	}

	if item.Ok {
		return nil
	}

	sql = `
	INSERT INTO core.KEYVALUES(_ID, VALUE, INDEX)
	VALUES ($1, $2, $3);`

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
	sql := `
	SELECT *
	FROM core.KEYVALUES
	WHERE _ID = $1
	LIMIT 1;`

	item, err := s.One(sql, key)
	if err != nil {
		return et.KeyValue{}, err
	}

	return et.KeyValue{
		Ok:    item.Ok,
		Value: item.Str("value"),
		Imdex: item.Int("index"),
	}, nil
}

/**
* DeleteKey - Delete key value
* @params key string
* @return error
**/
func (s *Postgres) DeleteKey(key string) error {
	sql := `
	DELETE 
	FROM core.KEYVALUES
	WHERE _ID = $1;`

	err := s.Exec(sql)
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

	result, err := s.One(sql, search)
	if err != nil {
		return et.List{}, err
	}

	all := result.Int("all")

	sql = `
	SELECT A.*
	FROM core.KEYVALUES A
	WHERE A.VALUE ILIKE %$1%
	OFFSET $2 LIMIT $3
	ORDER BY A.INDEX;`

	offset := (page - 1) * rows
	items, err := s.SQL(sql, search, offset, rows)
	if err != nil {
		return et.List{}, err
	}

	return items.ToList(all, page, rows), nil
}

package postgres

import (
	"database/sql"
	"time"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

func (s *Postgres) defineCache() error {
	exist, err := s.existTable("core", "CACHE")
	if err != nil {
		return console.Panic(err)
	}

	if exist {
		return s.defineRecordsFunction()
	}

	sql := parceSQL(`
  CREATE TABLE IF NOT EXISTS core.CACHE(		
		_ID VARCHAR(80) DEFAULT '',
		VALUE BYTEA,
		EXPIRATION TIMESTAMP DEFAULT NOW(),
		_IDT VARCHAR(80) DEFAULT '-1',
		INDEX BIGINT DEFAULT 0,
		PRIMARY KEY(_ID)
	);
	CREATE INDEX IF NOT EXISTS CACHE_EXPIRATION_IDX ON core.CACHE(EXPIRATION);
	CREATE INDEX IF NOT EXISTS CACHE__IDT_IDX ON core.CACHE(_IDT);
	CREATE INDEX IF NOT EXISTS CACHE_INDEX_IDX ON core.CACHE(INDEX);`)
	sql = strs.Append(sql, defineRecordTrigger("core.CACHE"), "\n")
	sql = strs.Append(sql, defineSeriesTrigger("core.CACHE"), "\n")

	err = s.Exec(sql)
	if err != nil {
		return console.Panic(err)
	}

	return s.defineRecordsFunction()
}

/**
* SetCache - Set key value
* @params key string, value []byte, duration time.Duration
* @return error
**/
func (s *Postgres) SetCache(key string, value []byte, duration time.Duration) error {
	sql := parceSQL(`
	UPDATE core.CACHE SET
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
	INSERT INTO core.CACHE(_ID, VALUE, INDEX, EXPIRATION)
	VALUES ($1, $2, $3, $4);`)

	expiration := time.Now().Add(duration)
	index := s.GetSerie("core.CACHE")
	err = s.Exec(sql, key, value, index, expiration)
	if err != nil {
		return err
	}

	return nil
}

/**
* GetCache - Get key value
* @params key string
* @return KeyValue, error
**/
func (s *Postgres) GetCache(key string) (et.KeyValue, error) {
	query := parceSQL(`
	SELECT VALUE, INDEX
	FROM core.CACHE
	WHERE _ID = $1
	LIMIT 1;`)

	var ok bool
	var value []byte
	var index int
	err := s.db.QueryRow(query, key).Scan(&value, &index)
	if err != nil {
		ok = err == sql.ErrNoRows
		if !ok {
			return et.KeyValue{}, err
		}
	}

	return et.KeyValue{
		Ok:    !ok,
		Value: value,
		Imdex: index,
	}, nil
}

/**
* DeleteCache - Delete key value
* @params key string
* @return error
**/
func (s *Postgres) DeleteCache(key string) error {
	sql := parceSQL(`
	DELETE 
	FROM core.CACHE
	WHERE _ID = $1;`)

	err := s.Exec(sql, key)
	if err != nil {
		return err
	}

	return nil
}

/**
* CleanCache - Clean cache
* @return error
**/
func (s *Postgres) CleanCache() error {
	sql := parceSQL(`
	DELETE
	FROM core.CACHE
	WHERE EXPIRATION < NOW();`)

	err := s.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

/**
* FindCache - Find key value
* @params search string, page, rows int
* @return et.List, error
**/
func (s *Postgres) FindCache(search string, page, rows int) (et.List, error) {
	sql := parceSQL(`
	SELECT COUNT(*) AS ALL
	FROM core.CACHE A
	WHERE A._ID ILIKE %$1%;`)

	result, err := s.Query(sql, search)
	if err != nil {
		return et.List{}, err
	}

	all := result.Int(0, "all")

	sql = parceSQL(`
	SELECT A.*
	FROM core.CACHE A
	WHERE A._ID ILIKE %$1%
	OFFSET $2 LIMIT $3
	ORDER BY A.INDEX;`)

	offset := (page - 1) * rows
	items, err := s.Query(sql, search, offset, rows)
	if err != nil {
		return et.List{}, err
	}

	return items.ToList(all, page, rows), nil
}

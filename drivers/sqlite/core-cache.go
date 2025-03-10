package sqlite

import (
	"time"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
)

func (s *SqlLite) defineCache() error {
	sql := parceSQL(`
  CREATE TABLE IF NOT EXISTS core.CACHE(		
		_ID VARCHAR(80) DEFAULT '',
		VALUE BYTEA,
		EXPIRATION TIMESTAMP DEFAULT NOW(),
		_IDT VARCHAR(80) DEFAULT '-1',
		INDEX BIGINT DEFAULT 0,
		PRIMARY KEY(_ID)
	);`)

	err := s.Exec(sql)
	if err != nil {
		return console.Panic(err)
	}

	return nil
}

func (s *SqlLite) SetCache(key string, value []byte, duration time.Duration) error

func (s *SqlLite) GetCache(key string) (et.KeyValue, error)
func (s *SqlLite) DeleteCache(key string) error
func (s *SqlLite) CleanCache() error
func (s *SqlLite) FindCache(search string, page, rows int) (et.List, error)

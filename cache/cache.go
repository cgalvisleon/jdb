package cache

import (
	"sync"
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/jdb/jdb"
)

type Cache struct {
	Id     string            `json:"id"`
	Values map[string]string `json:"_"`
	lock   *sync.RWMutex     `json:"-"`
	db     *jdb.DB           `json:"-"`
	source *jdb.Model        `json:"-"`
}

/**
* SetDB
* @param db *jdb.DB
**/
func (s *Cache) SetDB(db *jdb.DB) {
	s.db = db
}

/**
* Set
* @param key, value string, expiration time.Duration
* @return string
**/
func (s *Cache) Set(key, value string, expiration time.Duration) string {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Values[key] = value

	clean := func() {
		s.Delete(key)
	}

	duration := expiration * time.Second
	if duration != 0 {
		go time.AfterFunc(duration, clean)
	}

	event.Publish(
		"cache:set",
		et.Json{
			"now":        utility.NowTime(),
			"key":        key,
			"value":      value,
			"expiration": expiration,
			"originId":   s.Id,
		},
	)

	if s.source == nil {
		return value
	}

	go func() {
		exist, err := s.source.
			Where("key").Eq(key).
			Exist()
		if err != nil {
			return
		}

		if exist {
			_, err := s.source.Update(et.Json{
				"value": value,
			}).
				Where("key").Eq(key).
				Exec()
			if err != nil {
				return
			}
		} else {
			_, err := s.source.Insert(et.Json{
				"key":   key,
				"value": value,
			}).
				Exec()
			if err != nil {
				return
			}
		}
	}()

	return value
}

/**
* Get
* @param key string
* @return string, bool
**/
func (s *Cache) Get(key, def string) (string, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	result, ok := s.Values[key]
	if !ok && s.source != nil {
		item, err := s.source.
			Where("key").Eq(key).
			Select("value").
			One()
		if err != nil {
			return def, false
		}

		if !item.Ok {
			return def, false
		}

		result := item.Str("value")
		s.Set(key, result, 0)

		return result, true
	} else if !ok {
		return def, false
	}

	return result, true
}

/**
* Delete
* @param key string
**/
func (s *Cache) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.Values, key)

	if s.source == nil {
		return
	}

	go func() {
		event.Publish(
			"cache:delete",
			et.Json{
				"key":      key,
				"originId": s.Id,
			},
		)

		if s.source != nil {
			_, err := s.source.Delete().
				Where("key").Eq(key).
				Exec()
			if err != nil {
				return
			}
		}
	}()
}

var conn *Cache

func init() {
	if conn != nil {
		return
	}

	conn = &Cache{
		Id:     utility.RecordId("cache", ""),
		Values: make(map[string]string),
		lock:   &sync.RWMutex{},
	}
}

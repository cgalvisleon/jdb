package jdb

import (
	"sync"
	"time"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/utility"
)

type Cache struct {
	Id     string            `json:"id"`
	Values map[string][]byte `json:"_"`
	lock   *sync.RWMutex     `json:"-"`
	Db     *DB               `json:"-"`
}

/**
* Set
* @param key string, value []byte
* @return []byte
**/
func (s *Cache) Set(key string, value []byte, expiration time.Duration) []byte {
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

	return value
}

/**
* Get
* @param key string
* @return []byte
**/
func (s *Cache) Get(key string) []byte {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.Values[key]
}

/**
* Delete
* @param key string
**/
func (s *Cache) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.Values, key)
}

var cache *Cache

func init() {
	if cache != nil {
		return
	}

	cache = &Cache{
		Id:     utility.RecordId("cache", ""),
		Values: make(map[string][]byte),
		lock:   &sync.RWMutex{},
	}
}

func InitCacheEvents() {
	err := event.Subscribe("cache:set", eventCacheSet)
	if err != nil {
		console.Error(err)
	}

	err = event.Subscribe("cache:delete", eventCacheDelete)
	if err != nil {
		console.Error(err)
	}
}

func eventCacheSet(event event.EvenMessage) {
	data := event.Data
	key := data.Str("key")
	value := data.Str("value")
	originNow := data.Time("now")
	second := data.Num("expiration")
	originId := data.Str("originId")
	now := utility.NowTime()
	diference := now.Sub(originNow)
	expiration := time.Duration(second) - diference

	if originId != cache.Id {
		cache.Set(key, []byte(value), expiration)
		console.Logf("Cache set", `Key:%s Value:%s Expirate:%v`, key, value, expiration)
	}
}

func eventCacheDelete(event event.EvenMessage) {
	data := event.Data
	key := data.Str("key")
	originId := data.Str("originId")

	if originId != cache.Id {
		cache.Delete(key)
		console.Logf("Cache delete", `Key:%s`, key)
	}
}

/**
* SetCache - Set key value
* @params key string, value string, expiration time.Duration
* @return error
**/
func SetCache(key string, value string, expiration time.Duration) {
	now := utility.NowTime()
	cache.Set(key, []byte(value), expiration)
	event.Publish("cache:set", et.Json{"key": key, "value": value, "now": now, "expiration": expiration, "originId": cache.Id})

	if cache.Db != nil {
		go cache.Db.SetCache(key, []byte(value), expiration)
	}
}

/**
* GetCache - Get key value
* @params key string
* @return et.KeyValue
**/
func GetCache(key string) et.KeyValue {
	value := cache.Get(key)
	if value == nil && cache.Db != nil {
		kv, err := cache.Db.GetCache(key)
		if err != nil {
			return et.KeyValue{}
		}

		value = kv.Value
		cache.Set(key, value, 0)

		return kv
	}

	return et.KeyValue{}
}

/**
* DeleteCache - Delete key value
* @params key string
* @return error
**/
func DeleteCache(key string) {
	cache.Delete(key)
	event.Publish("cache:delete", et.Json{"key": key, "originId": cache.Id})

	if cache.Db != nil {
		go cache.Db.DeleteCache(key)
	}
}

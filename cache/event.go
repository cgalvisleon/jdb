package cache

import (
	"time"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/utility"
)

func InitEvents() {
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

	if originId != conn.Id {
		conn.Set(key, value, expiration)
		console.Logf("Cache set", `Key:%s Value:%s Expirate:%v`, key, value, expiration)
	}
}

func eventCacheDelete(event event.EvenMessage) {
	data := event.Data
	key := data.Str("key")
	originId := data.Str("originId")

	if originId != conn.Id {
		conn.Delete(key)
		console.Logf("Cache delete", `Key:%s`, key)
	}
}

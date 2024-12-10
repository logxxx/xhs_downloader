package utils

import (
	cache "github.com/patrickmn/go-cache"
	"time"
)

var (
	_cache = cache.New(10*time.Minute, 30*time.Minute)
)

func TryOccupyKey(key string, dura time.Duration, occupySuccFn func(), occupyFailFn ...func()) {
	if key == "" {
		return
	}
	if dura.Seconds() < 1 {
		dura = 1 * time.Second
	}
	_, exists := _cache.Get(key)
	if exists {
		if len(occupyFailFn) > 0 && occupyFailFn[0] != nil {
			occupyFailFn[0]()
		}
		return
	}
	if occupySuccFn != nil {
		occupySuccFn()
	}
	_cache.Set(key, time.Now().Unix(), dura)
	return
}

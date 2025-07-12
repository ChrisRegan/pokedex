package pokecache

import (
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type cache struct {
	entries map[string]cacheEntry
}

func NewCache(interval time.Duration) {
	cache := cache{}
	cache.reapLoop(interval)
}

func (cache cache) Add(key string, val []byte) {

	cache.entries[key] = cacheEntry{createdAt: time.Now(), val: val}

}

func (cache cache) Get(key string) (val cacheEntry, found bool) {
	val, found = cache.entries[key]
	return val, found
}

func (cache cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		for k := range cache.entries {
			delete(cache.entries, k)
		}
	}
}

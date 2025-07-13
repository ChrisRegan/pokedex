package pokecache

import (
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	shutdown chan bool
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{entries: make(map[string]cacheEntry), shutdown: make(chan bool)}
	cache.reapLoop(interval)
	return cache
}

func (cache Cache) ShutdownCacheCollection() {
	cache.shutdown <- true
}

func (cache Cache) Add(key string, val []byte) {

	cache.entries[key] = cacheEntry{createdAt: time.Now(), val: val}

}

func (cache Cache) Get(key string) ([]byte, bool) {
	entry, found := cache.entries[key]
	return entry.val, found
}

func (cache Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-cache.shutdown:
				cache.entries = nil
				ticker.Stop()
				return
			case <-ticker.C:
				for k := range cache.entries {
					delete(cache.entries, k)
				}
			}
		}
	}()

}

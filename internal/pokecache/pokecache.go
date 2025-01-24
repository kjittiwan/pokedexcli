package pokecache

import (
	"sync"
	"time"
)

// entry's key will be the url called, val will be the raw data
type Cache struct {
	mu    sync.Mutex
	entry map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// create new cache and return
func NewCache(interval time.Duration) *Cache {
	// create ticker for duration of interval
	ticker := time.NewTicker(interval)
	// when interval pass, delete old entries
	cache := Cache{
		entry: map[string]cacheEntry{},
	}
	go func(ch <-chan time.Time) {
		for item := range ch {
			t := item
			// delete old
			cache.readLoop(t.Add(-interval))
		}
	}(ticker.C)
	return &cache
}

// add a new entry to the cache
func (c *Cache) Add(key string, val []byte) {
	// adds new entry to the cache
	c.mu.Lock()
	c.entry[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mu.Unlock()
}

// get value of cache from given key, bool represents success
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	val, ok := c.entry[key]
	c.mu.Unlock()
	if !ok {
		return []byte{}, false
	}
	return val.val, true
}

// called when cache is created, remove caches older than a certain point
func (c *Cache) readLoop(t time.Time) {
	for key, val := range c.entry {
		if val.createdAt.Before(t) {
			c.mu.Lock()
			delete(c.entry, key)
			c.mu.Unlock()
		}
	}
}

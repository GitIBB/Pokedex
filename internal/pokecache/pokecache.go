package pokecache

import (
	"sync"
	"time"
)

// entries maps keys to cached values with creation timestamps
// mu ensured thread-safe access to entire map
type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

// cacheEntry represents a single item in the cache with its metadata
type cacheEntry struct {
	createdAt time.Time // timestamp of entry creation
	val       []byte    // cached data
}

// creates a Cache isntance and commences background cleanup routine by interval
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries: make(map[string]cacheEntry),
	}

	// reaping loop for periodically removing expired entries
	go cache.reapLoop(interval)

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// create a new cacheEntry
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	//add cacheEntry to map
	c.entries[key] = entry

}

// Get retrieves value from cache by key
// returns value and boolean indicating if key found
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()         //mutex lock
	defer c.mu.Unlock() // unlock deferral on return

	entry, exists := c.entries[key] //attempt to get entry from cache
	if exists {
		return entry.val, true // if found, return value and true
	}
	return nil, false
}

// reapLoop runs on a timer to clean expired cache entries
// interval determines how often checks are done and maximum entry age
func (c *Cache) reapLoop(interval time.Duration) {

	ticker := time.NewTicker(interval) // create interval ticker

	defer ticker.Stop() // ticker is cleaned on return

	// infinite for loop to continuously check for expired entries
	for {
		<-ticker.C // wait for next tick

		c.mu.Lock() // lock cache with mutex to prevent concurrent access while cleaning

		// check cache entries
		for key, entry := range c.entries {
			//if entry is older than interval, delete
			if time.Since(entry.createdAt) > interval {
				delete(c.entries, key)
			}
		}
		//unlock cache
		c.mu.Unlock()
	}

}

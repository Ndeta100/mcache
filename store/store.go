package store

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Ndeta100/mcache/config"
)

type Item struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	store          map[string]Item
	mutex          sync.RWMutex
	ttl            time.Duration
	capacity       int
	stopClean      chan struct{}
	filename       string
	backupfile     string
	snapshotchange bool
}

func NewCache(options config.CacheOptions) *Cache {
	cache := &Cache{
		store:      make(map[string]Item, options.Capacity),
		ttl:        options.TTL,
		capacity:   options.Capacity,
		stopClean:  make(chan struct{}),
		filename:   "cache_snap.gob",
		backupfile: "cache_snap_backup.gob",
	}
	// Load the cache from the YAML snapshot file
	cache.LoadFromDisk()
	//schedule eviction
	go cache.scheduleCacheEvict(3)
	// Start a background Goroutine to periodically save snapshots
	go cache.startSnapshotRoutine()
	return cache
}

// Set Insert or update values
func (c *Cache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	key = strings.ToLower(key)
	//Check if element exist before adding
	found := c.isValidEntry(key, value)
	if found {
		fmt.Printf("Element already exist %s\n", value)
		return
	}
	// Strictly enforce capacity limit: If at capacity, do not add new items
	if len(c.store) >= c.capacity {
		fmt.Printf("Cache capacity %d reached, cannot add new item with key '%s'\n", c.capacity, key)
		return
	}
	// Normalize the key to ensure consistency
	key = strings.ToLower(key)
	// Calculate expiration time
	var expiration int64
	if c.ttl > 0 {
		expiration = c.getExpiration()
	} else {
		expiration = 0
	}
	// If the item already exists, update its value and expiration time
	if item, found := c.store[key]; found {
		item.Value = value
		item.Expiration = int64(expiration)
		c.store[key] = item
		return
	}
	// Add new item to the cache
	c.store[key] = Item{
		Value:      value,
		Expiration: int64(expiration),
	}
	c.snapshotchange = true
	c.LogOperation(key, value)

}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// Normalize the key to ensure consistency
	key = strings.ToLower(key)
	item, found := c.store[key]
	if !found {
		return nil, false
	}
	if item.Expiration > 0 && time.Now().Unix() > int64(item.Expiration) {
		return "value expire", false
	}

	return item.Value, true
}

func (c *Cache) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// Normalize the key to ensure consistency
	key = strings.ToLower(key)
	_, found := c.store[key]
	if !found {
		fmt.Printf("item with key %s not found", key)
	}
	delete(c.store, key)
	c.snapshotchange = true
	return true
}

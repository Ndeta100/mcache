package store

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Ndeta100/config"
)

type Item struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	store     map[string]Item
	mutex     sync.RWMutex
	ttl       time.Duration
	capacity  int
	stopClean chan struct{}
}

func NewCache(options config.CacheOptions) *Cache {
	return &Cache{
		store:     make(map[string]Item, options.Capacity),
		ttl:       options.TTL,
		capacity:  options.Capacity,
		stopClean: make(chan struct{}),
	}
}

// Insert or update values
func (c *Cache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// Strictly enforce capacity limit: If at capacity, do not add new items
	if len(c.store) >= c.capacity {
		fmt.Printf("Cache capacity %d reached, cannot add new item with key '%s'\n", c.capacity, key)
		return
	}
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
		Expiration: int64(expiration), // Added trailing comma for proper formatting
	}

}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item, found := c.store[key]
	if !found {
		return nil, false
	}
	if item.Expiration > 0 && time.Now().Unix() > int64(item.Expiration) {
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) backgroundClean() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			c.cleanUp()
		case <-c.stopClean:
			c.StopCleanup()
		}
	}
}

func (c *Cache) cleanUp() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	now := time.Now().Unix()
	for key, item := range c.store {
		if item.Expiration > 0 && now > item.Expiration {
			delete(c.store, key)
		}
	}
}

// get expirey date
func (c *Cache) getExpiration() int64 {
	if c.ttl <= 0 {
		return 0
	}
	return time.Now().Add(c.ttl).Unix()
}

// StopCleanup stops the background cleanup goroutine.
func (c *Cache) StopCleanup() {
	close(c.stopClean)
}

func (c *Cache) String() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var b strings.Builder
	fmt.Fprintf(&b, "Cache{capacity: %d, items: %d, ", c.capacity, len(c.store))

	if c.ttl > 0 {
		fmt.Fprintf(&b, "ttl: %v, ", c.ttl)
	}

	b.WriteString("items: [")
	for k, v := range c.store {
		fmt.Fprintf(&b, "{%s: %v (expires: %v)}, ", k, v.Value, time.Unix(v.Expiration, 0))
	}
	b.WriteString("]}")

	return b.String()
}

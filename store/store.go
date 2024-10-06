package store

import (
	"encoding/gob"
	"fmt"
	"os"
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

	// Start a background Goroutine to periodically save snapshots
	go cache.startSnapshotRoutine()
	return cache
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
	c.snapshotchange = true
	c.LogOperation(key, value)

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

// SaveToDisk saves the current cache state to a binary file.
func (c *Cache) SaveToDisk() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// If no changes were made, there's no need to save the cache
	if !c.snapshotchange {
		return
	}

	// Open the file for writing. If the file does not exist, it will be created.
	// If it does exist, it will be truncated and overwritten.
	temp := c.filename + ".temp"
	file, err := os.OpenFile(temp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening snapshot file: %v\n", err)
		return
	}
	defer file.Close()

	// Encode the entire cache to the file
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(c.store); err != nil {
		fmt.Printf("Error encoding cache to gob: %v\n", err)
		return
	}
	if err := file.Sync(); err != nil {
		fmt.Errorf("error syncing file: %v", err)
		return
	}
	if err := os.Rename(temp, c.filename); err != nil {
		fmt.Errorf("error renaming temp file: %v", err)
		return
	}
	fmt.Println("Cache snapshot saved to disk successfully.")
	c.snapshotchange = false
}

// LoadFromDisk loads the cache state from a binary file.
func (c *Cache) LoadFromDisk() {
	file, err := os.Open(c.filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No snapshot found, starting with an empty cache.")
			return
		}
		fmt.Printf("Error opening snapshot file: %v\n", err)
		return
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var store map[string]Item
	if err := decoder.Decode(&store); err != nil {
		fmt.Printf("Error decoding gob cache data: %v\n", err)
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.store = store
	fmt.Println("Cache state loaded from disk successfully.")
}

// SaveToDiskAsync saves the cache state asynchronously to avoid blocking.
func (c *Cache) SaveToDiskAsync() {
	go func() {
		c.SaveToDisk()
	}()
}

// LogOperation logs a cache write operation to an append-only log file.
func (c *Cache) LogOperation(key string, value interface{}) {
	file, err := os.OpenFile("cache_aof.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening AOF file: %v\n", err)
		return
	}
	defer file.Close()

	// Record operation in the form of "SET key value"
	logLine := fmt.Sprintf("SET %s %v\n", key, value)
	if _, err := file.WriteString(logLine); err != nil {
		fmt.Printf("Error writing to AOF: %v\n", err)
	}
}

// startSnapshotRoutine starts a Goroutine to periodically save the cache to disk.
func (c *Cache) startSnapshotRoutine() {
	ticker := time.NewTicker(1 * time.Minute) // Snapshot interval: 1 minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Save the snapshot asynchronously to avoid blocking the main routine
			c.SaveToDiskAsync()
		case <-c.stopClean:
			return
		}
	}
}

// StopSnapshotRoutine stops the periodic snapshot routine.
func (c *Cache) StopSnapshotRoutine() {
	close(c.stopClean)
}

// get expirey date
func (c *Cache) getExpiration() int64 {
	if c.ttl <= 0 {
		return 0
	}
	return time.Now().Add(c.ttl).Unix()
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

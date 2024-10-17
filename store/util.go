package store

import (
	"encoding/gob"
	"fmt"
	"os"
	"strings"
	"time"
)

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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file", err)
		}
	}(file)

	// Encode the entire cache to the file
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(c.store); err != nil {
		fmt.Printf("Error encoding cache to gob: %v\n", err)
		return
	}
	if err := file.Sync(); err != nil {
		_ = fmt.Errorf("error syncing file: %v", err)
		return
	}
	if err := os.Rename(temp, c.filename); err != nil {
		_ = fmt.Errorf("error renaming temp file: %v", err)
		return
	}
	fmt.Println("Cache snapshot saved to disk successfully.")
	c.snapshotchange = false
}

// SaveToDiskAsync saves the cache state asynchronously to avoid blocking.
func (c *Cache) SaveToDiskAsync() {
	go func() {
		c.SaveToDisk()
	}()
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
	now := time.Now().Unix()
	fileStore := make(map[string]Item)
	for key, item := range store {
		if item.Expiration > 0 && item.Expiration < now {
			delete(c.store, key)
			fmt.Printf("Cache removed expired item: %v\n", key)
		}
	}
	c.store = fileStore
	fmt.Println("Cache state loaded from disk successfully.")
}

// startSnapshotRoutine starts a Goroutine to periodically save the cache to disk.
func (c *Cache) startSnapshotRoutine() {
	ticker := time.NewTicker(1 * time.Second) // Snapshot interval: 1 minute
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

// get expire date
func (c *Cache) getExpiration() int64 {
	if c.ttl <= 0 {
		return 0
	}
	return time.Now().Add(c.ttl).Unix()
}

func (c *Cache) isValidEntry(key string, value interface{}) bool {
	element, exists := c.store[key]
	if !exists {
		return false
	}

	// Check if the value matches and the item hasn't expired
	if element.Value == value && (element.Expiration == 0 || element.Expiration > time.Now().Unix()) {
		return true
	}
	return false
}

func (c *Cache) scheduleEviction(key string) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	select {
	case <-c.stopClean:

	}
}
func (c *Cache) evictExpiredItems() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	now := time.Now().Unix()
	for key, item := range c.store {
		if item.Expiration > 0 && item.Expiration <= now {
			delete(c.store, key)
			fmt.Printf("Evicting expired item:%v\n", key)
		}
	}
}

func (c *Cache) scheduleCacheEvict(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			c.evictExpiredItems()
		}
	}()
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

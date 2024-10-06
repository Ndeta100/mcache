package store

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Ndeta100/config"
)

// Test for Set and Get operations
func TestCacheSetAndGet(t *testing.T) {
	// Create a cache with TTL of 10 seconds and capacity of 10
	cfgCache := config.CacheOptions{
		TTL:      time.Second * 10,
		Capacity: 10,
	}
	cache := NewCache(cfgCache)

	// Set a key-value pair in the cache
	cache.Set("key1", "value1")

	// Get the value from the cache
	value, found := cache.Get("key1")
	if !found {
		t.Errorf("Expected to find key1, but it was not found")
	}
	if value != "value1" {
		t.Errorf("Expected value 'value1', got %v", value)
	}
}

// Test for expiration
func TestCacheExpiration(t *testing.T) {
	// Create a cache with TTL of 1 second and capacity of 10
	cfgCache := config.CacheOptions{
		TTL:      time.Second * 1,
		Capacity: 10,
	}
	cache := NewCache(cfgCache)

	// Set a key-value pair in the cache
	cache.Set("key2", "value2")

	// Wait for 2 seconds to let the item expire
	time.Sleep(2 * time.Second)

	// Get the value from the cache
	_, found := cache.Get("key2")
	if found {
		t.Errorf("Expected key2 to be expired, but it was found")
	}
}

// Test updating an existing item
func TestCacheUpdate(t *testing.T) {
	// Create a cache with TTL of 10 seconds and capacity of 10
	cfgCache := config.CacheOptions{
		TTL:      time.Second * 10,
		Capacity: 10,
	}
	cache := NewCache(cfgCache)

	// Set a key-value pair in the cache
	cache.Set("key3", "value3")

	// Update the value for the same key
	cache.Set("key3", "newValue3")

	// Get the updated value from the cache
	value, found := cache.Get("key3")
	if !found {
		t.Errorf("Expected to find key3, but it was not found")
	}
	if value != "newValue3" {
		t.Errorf("Expected value 'newValue3', got %v", value)
	}
}

// Test for cache capacity and eviction logic
func TestCacheCapacity(t *testing.T) {
	// Create a cache with TTL of 10 seconds and capacity of 2
	cfgCache := config.CacheOptions{
		TTL:      time.Second * 10,
		Capacity: 2,
	}
	cache := NewCache(cfgCache)

	// Set items in the cache to exceed capacity
	cache.Set("key4", "value4")
	cache.Set("key5", "value5")
	cache.Set("key6", "value6") // This should trigger eviction

	// One of the previous keys should have been evicted
	_, found4 := cache.Get("key4")
	_, found5 := cache.Get("key5")
	_, found6 := cache.Get("key6")

	// Check if at least one of the earlier keys has been evicted to accommodate key6
	if !found4 && !found5 && !found6 {
		t.Errorf("Expected one of the keys to be in the cache, but none were found")
	}
}

// Test for concurrency safety
func TestCacheConcurrency(t *testing.T) {
	// Create a cache with TTL of 10 seconds and capacity of 10
	cfgCache := config.CacheOptions{
		TTL:      time.Second * 10,
		Capacity: 10,
	}
	cache := NewCache(cfgCache)

	// Create a WaitGroup
	var wg sync.WaitGroup

	// Set the number of goroutines we're going to run
	wg.Add(10)

	// Run multiple goroutines to access the cache concurrently
	for i := 0; i < 10; i++ {
		go func(n int) {
			// Ensure the WaitGroup is decremented when the goroutine completes
			defer wg.Done()

			key := "key" + strconv.Itoa(n)
			cache.Set(key, n)
			value, found := cache.Get(key)

			// Add some assertions
			if !found {
				t.Errorf("Key %s was not found in the cache", key)
			}
			if value != n {
				t.Errorf("Expected value %d for key %s, but got %v", n, key, value)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
}

package store

import (
	"sync"
	"time"

	"github.com/Ndeta100/config"
)

type Item struct {
	Value      interface{}
	Expiration int16
}

type Cache struct {
	store    map[string]Item
	mutex    sync.RWMutex
	ttl      time.Duration
	capacity int
}

func NewCache(options config.CacheOptions) *Cache {
	return &Cache{
		store:    make(map[string]Item, options.Capacity),
		ttl:      options.TTL,
		capacity: options.Capacity,
	}
}

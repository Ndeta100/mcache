package main

import (
	"fmt"
	"time"

	"github.com/Ndeta100/config"
	"github.com/Ndeta100/server"
	"github.com/Ndeta100/store"
)

func main() {
	cfgcache := config.CacheOptions{
		TTL:      time.Second * 10,
		Capacity: 10,
	}
	cache := store.NewCache(cfgcache)

	// for i := 0; i < 12; i++ {
	// 	cache.Set(fmt.Sprintf("key_%v", i), fmt.Sprintf("Banana_%v", i))
	// }
	// Get the value from the cache
	// value, found := cache.Get("keyq")
	cache.SaveToDisk()
	fmt.Println(cache.String())
	// if found {
	// 	fmt.Println("Found value:", value)
	// } else {
	// 	fmt.Println("Value not found or expired")
	// }

	server := server.NewServer()
	server.Start()

}

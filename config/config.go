package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Cache settings.
	Cache CacheOptions `yaml:"cache"`

	// Network settings.
	Network NetworkSettings `yaml:"network"`

	// Logging settings.
	Logging LoggingOptions `yaml:"logging"`

	// Security settings.
	Security SecurityOptions `yaml:"security"`

	// Storage settings (optional, if persistence is needed).
	Storage StorageOptions `yaml:"storage,omitempty"`
}

// CacheOptions for the cache module.
type CacheOptions struct {
	// Maximum number of items in the cache.
	Capacity int           `yaml:"capacity" default:"1000"`
	TTL      time.Duration `yaml:"ttl"` // Time to live for each item.

	// Eviction policy strategy (e.g., LRU, FIFO, LFU).
	Strategy string `yaml:"strategy" default:"LRU"`
}

// NetworkSettings Network settings for the network module.
type NetworkSettings struct {
	Host string `yaml:"host" default:"127.0.0.1"`
	Port int    `yaml:"port" default:"6379"`
}

// LoggingOptions Logging settings for the logging module.
type LoggingOptions struct {
	Level  string `yaml:"level" default:"info"`
	Format string `yaml:"format" default:"text"`
}

// SecurityOptions Security options for the security module.
type SecurityOptions struct {
	// Optional authentication settings.
	RequireAuth bool   `yaml:"require_auth" default:"false"`
	Password    string `yaml:"password"` // Only if RequireAuth is true.
}

// StorageOptions Storage options for persistence (optional).
type StorageOptions struct {
	EnablePersistence bool   `yaml:"enable_persistence" default:"false"`
	StorageFile       string `yaml:"storage_file" default:"cache_data.dat"`
}

// LoadConfiguration loads the configuration from a YAML file.
func LoadConfiguration(cfg *Config, filename string) error {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// if the LoadConfiguration does not exist we write one to the filename
		*cfg = defaultConfig
		data, err := yaml.Marshal(cfg)
		if err != nil {
			fmt.Println("Error with yaml Marshal", err)
		}
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			return err
		}
	} else {
		_, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("Error with yaml Marshal", err)
		}
	}

	// //TODO: Validate cfg if needed
	// if err := cfg.Validate(); err != nil {
	// 	return err
	// }

	return nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	// Validate network port
	if c.Network.Port <= 0 || c.Network.Port > 65535 {
		return fmt.Errorf("invalid network port: %d", c.Network.Port)
	}

	// Validate cache strategy
	validStrategies := map[string]bool{
		"LRU":  true,
		"FIFO": true,
		"LFU":  true,
	}
	if !validStrategies[c.Cache.Strategy] {
		return fmt.Errorf("invalid cache strategy: %s", c.Cache.Strategy)
	}

	//TODO: Additional validations

	return nil
}

func GetDefaultConfig() Config {
	return defaultConfig
}

// Default configuration options.
var defaultConfig = Config{
	Cache: CacheOptions{
		Capacity: 1000,
		TTL:      0, // 0 means no expiration
		Strategy: "LRU",
	},
	Network: NetworkSettings{
		Host: "", //"127.0.0.1"
		Port: 6379,
	},
	Logging: LoggingOptions{
		Level:  "info",
		Format: "text",
	},
	Security: SecurityOptions{
		RequireAuth: false,
		Password:    "",
	},
	Storage: StorageOptions{
		EnablePersistence: false,
		StorageFile:       "cache_data.dat",
	},
}

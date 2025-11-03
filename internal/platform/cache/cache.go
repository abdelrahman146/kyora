// Package cache provides caching utilities
// implement memcached
package cache

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/spf13/viper"
)

type Cache struct {
	mc *memcache.Client
}

func NewConnection() *Cache {
	servers := viper.GetStringSlice(config.CacheHosts)
	mc := memcache.New(servers...)
	return &Cache{mc: mc}
}

// Get retrieves a value from the cache.
func (m *Cache) Get(key string) ([]byte, error) {
	item, err := m.mc.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

// SetX sets a value in the cache with an expiration time.
func (m *Cache) SetX(key string, value []byte, expiration int32) error {
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}
	return m.mc.Set(item)
}

// Set sets a value in the cache without expiration.
func (m *Cache) Set(key string, value []byte) error {
	item := &memcache.Item{
		Key:   key,
		Value: value,
	}
	return m.mc.Set(item)
}

// Delete removes a value from the cache.
func (m *Cache) Delete(key string) error {
	return m.mc.Delete(key)
}

func (m *Cache) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (m *Cache) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Increment atomically increments a numeric value in the cache.
// It initializes the key to 1 if it doesn't exist.
func (m *Cache) Increment(key string, delta uint64) error {
	_, err := m.mc.Increment(key, delta)
	if errors.Is(err, memcache.ErrCacheMiss) {
		// If the key doesn't exist, initialize it with the delta value.
		// Memcache's Increment starts from 0, so we set it to delta (usually 1).
		item := &memcache.Item{
			Key:   key,
			Value: []byte(strconv.FormatUint(delta, 10)),
		}
		return m.mc.Set(item)
	}
	return err
}

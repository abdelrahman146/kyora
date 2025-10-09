package db

import "github.com/bradfitz/gomemcache/memcache"

type Memcache struct {
	mc *memcache.Client
}

func NewMemcache(servers []string) *Memcache {
	mc := memcache.New(servers...)
	return &Memcache{mc: mc}
}

func (m *Memcache) Get(key string) ([]byte, error) {
	item, err := m.mc.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func (m *Memcache) Set(key string, value []byte, expiration int32) error {
	item := &memcache.Item{
		Key:   key,
		Value: value,
	}
	return m.mc.Set(item)
}

func (m *Memcache) Delete(key string) error {
	return m.mc.Delete(key)
}

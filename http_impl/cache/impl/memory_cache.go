package impl

import (
	"distributed-cache/cache"
	"fmt"
	"sync"
)

type MemoryCache struct {
	cache   map[string][]byte
	rwMutex sync.RWMutex
	cache.Stat
}

func (m *MemoryCache) Set(k string, v []byte) {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	if _, ok := m.cache[k]; ok {
		m.Stat.Del(k, v)
	}
	m.cache[k] = v
	m.Stat.Add(k, v)
}

func (m *MemoryCache) Get(k string) ([]byte, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	if v, ok := m.cache[k]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("this key not exists")
}

func (m *MemoryCache) Del(k string) error {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	if v, ok := m.cache[k]; ok {
		m.Stat.Del(k, v)
		delete(m.cache, k)
		return nil
	}
	return fmt.Errorf("this key not exists")
}

func (m *MemoryCache) GetStat() cache.Stat {
	return m.Stat
}


func NewMemoryCache()*MemoryCache {
	return &MemoryCache{cache: make(map[string][]byte), rwMutex: sync.RWMutex{}, Stat: cache.Stat{}}
}
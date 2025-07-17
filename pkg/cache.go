package pkg

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value     float64
	Timestamp time.Time
}

type RateCache struct {
	data  map[string]CacheItem
	mutex sync.RWMutex
	ttl   time.Duration
}

func NewRateCache(ttl time.Duration) *RateCache {
	return &RateCache{
		data:  make(map[string]CacheItem),
		ttl:   ttl,
		mutex: sync.RWMutex{},
	}
}

func (c *RateCache) Get(key string) (float64, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, found := c.data[key]
	if !found || time.Since(item.Timestamp) > c.ttl {
		return 0, false
	}
	return item.Value, true
}

func (c *RateCache) Set(key string, value float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = CacheItem{
		Value:     value,
		Timestamp: time.Now(),
	}
}

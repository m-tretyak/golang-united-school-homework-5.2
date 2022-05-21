package cache

import (
	"sync"
	"time"
)

type Cache struct {
	items  map[string]cacheItem
	locker sync.Locker
}

type cacheItem struct {
	value    string
	deadline time.Time
}

func (item *cacheItem) isExpired(now time.Time) bool {
	return item != nil && !item.deadline.IsZero() && item.deadline.Before(now)
}

func NewCache() (result *Cache) {
	result = new(Cache)
	result.locker = new(sync.Mutex)
	result.locker.Lock()
	defer result.locker.Unlock()

	result.items = make(map[string]cacheItem)
	return
}

func (cache *Cache) Get(key string) (string, bool) {
	cache.locker.Lock()
	defer cache.locker.Unlock()

	now := time.Now()

	if item, ok := cache.items[key]; ok && !item.isExpired(now) {
		return item.value, true
	} else if ok {
		delete(cache.items, key)
	}

	return "", false
}

func (cache *Cache) Put(key, value string) {
	cache.PutTill(key, value, time.Time{})
}

func (cache *Cache) Keys() (result []string) {
	cache.locker.Lock()
	defer cache.locker.Unlock()

	result = make([]string, 0, len(cache.items))
	now := time.Now()

	for key, item := range cache.items {
		if item.isExpired(now) {
			delete(cache.items, key)
			continue
		}

		result = append(result, key)
	}

	return
}

func (cache *Cache) PutTill(key, value string, deadline time.Time) {
	cache.locker.Lock()
	defer cache.locker.Unlock()

	cache.items[key] = cacheItem{
		value:    value,
		deadline: deadline,
	}
}

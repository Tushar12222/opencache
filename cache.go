package main

import "sync"

type Cache struct {
	kv    map[string]any
	mutex sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		kv: make(map[string]any),
	}
}

func (cache *Cache) keyExists(key string) bool {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	_, exists := cache.kv[key]
	return exists
}

func (cache *Cache) getKey(key string) (any, error) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	val, exists := cache.kv[key]
	if !exists {
		return nil, &NotFoundError{Message: "key not found"}
	}

	return val, nil
}

func (cache *Cache) putKey(key string, value any) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	_, exists := cache.kv[key]
	if !exists {
		return &NotFoundError{Message: "key not found"}
	}
	cache.kv[key] = value
	return nil
}

func (cache *Cache) deleteKey(key string) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	_, exists := cache.kv[key]
	if !exists {
		return &NotFoundError{Message: "key not found"}
	}

	delete(cache.kv, key)
	return nil
}

func (cache *Cache) addKey(key string, value any) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	_, exists := cache.kv[key]
	if exists {
		return &ConflictError{Message: "key already exists"}
	}

	cache.kv[key] = value
	return nil
}

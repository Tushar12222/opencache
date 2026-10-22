package main

import (
	"sync"
	"time"
)

type Value struct {
	val        any
	expiration *time.Time
}

type Cache struct {
	kv    map[string]Value
	mutex sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		kv: make(map[string]Value),
	}
}

func (cache *Cache) keyExists(key string) bool {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	_, exists := cache.kv[key]
	return exists
}

func (cache *Cache) getKey(key string) (any, *time.Time, error) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	val, exists := cache.kv[key]
	if !exists {
		return nil, nil, &NotFoundError{Message: "key not found"}
	}

	return val.val, val.expiration, nil
}

func (cache *Cache) putKey(key string, value any, expiration *time.Time) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	_, exists := cache.kv[key]
	if !exists {
		return &NotFoundError{Message: "key not found"}
	}
	if value != nil && expiration != nil {
		cache.kv[key] = Value{value, expiration}
	} else if value != nil {
		v := cache.kv[key]
		cache.kv[key] = Value{val: value, expiration: v.expiration}
	} else {
		v := cache.kv[key]
		cache.kv[key] = Value{val: v.val, expiration: nil}
	}
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

func (cache *Cache) addKey(key string, value any, expiration *time.Time) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	_, exists := cache.kv[key]
	if exists {
		return &ConflictError{Message: "key already exists"}
	}

	cache.kv[key] = Value{value, expiration}
	return nil
}

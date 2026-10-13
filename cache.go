package main

import "errors"

type Cache struct {
	kv map[any]any
}

func (cache *Cache) keyExists(key any) bool {
	_, exists := cache.kv[key]
	if !exists {
		return false
	}

	return true
}

func (cache *Cache) getKey(key any) (any, error) {
	val, exists := cache.kv[key]
	if !exists {
		return nil, errors.New("Key not found")
	}

	return val, nil
}

func (cache *Cache) putKey(key any, value any) error {
	_, exists := cache.kv[key]
	if !exists {
		return errors.New("Key not found")
	}
	cache.kv[key] = value
	return nil
}

func (cache *Cache) deleteKey(key any) error {
	_, exists := cache.kv[key]
	if !exists {
		return errors.New("Key not found")
	}

	delete(cache.kv, key)
	return nil
}

func (cache *Cache) addKey(key any, value any) error {
	_, exists := cache.kv[key]
	if exists {
		return errors.New("Key already exists")
	}

	cache.kv[key] = value
	return nil
}

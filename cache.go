package main

type Cache struct {
	kv map[any]any
}

func NewCache() *Cache {
	return &Cache{
		kv: make(map[any]any),
	}
}

func (cache *Cache) keyExists(key any) bool {
	_, exists := cache.kv[key]
	return exists
}

func (cache *Cache) getKey(key any) (any, error) {
	val, exists := cache.kv[key]
	if !exists {
		return nil, &NotFoundError{Message: "key not found"}
	}

	return val, nil
}

func (cache *Cache) putKey(key any, value any) error {
	_, exists := cache.kv[key]
	if !exists {
		return &NotFoundError{Message: "key not found"}
	}
	cache.kv[key] = value
	return nil
}

func (cache *Cache) deleteKey(key any) error {
	_, exists := cache.kv[key]
	if !exists {
		return &NotFoundError{Message: "key not found"}
	}

	delete(cache.kv, key)
	return nil
}

func (cache *Cache) addKey(key any, value any) error {
	_, exists := cache.kv[key]
	if exists {
		return &ConflictError{Message: "key already exists"}
	}

	cache.kv[key] = value
	return nil
}

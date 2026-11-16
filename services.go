package main

import "time"

type Service struct {
	cache *Cache
}

func NewService(cache *Cache) *Service {
	return &Service{
		cache: cache,
	}
}

func (service *Service) keyExists(key string) bool {
	now := time.Now().UTC()
	return service.cache.keyExists(key, now)
}

func (service *Service) getKey(key string) (byte, []byte, *time.Time, error) {
	now := time.Now().UTC()
	return service.cache.getKey(key, now)
}

func (service *Service) addKey(key string, valueType byte, valueBytes []byte, exp int) error {
	var expiration *time.Time
	var now time.Time
	if exp != 0 {
		now = time.Now().UTC()
		expirationTime := now.Add(time.Duration(exp) * time.Millisecond)
		expiration = &expirationTime
	}

	if now.IsZero() {
		now = time.Now().UTC()
	}
	return service.cache.addKey(key, valueType, valueBytes, now, expiration)
}

func (service *Service) putKey(key string, valueType byte, valueBytes []byte, exp int) error {
	var expiration *time.Time
	var now time.Time
	if exp != 0 {
		now := time.Now().UTC()
		expirationTime := now.Add(time.Duration(exp) * time.Millisecond)
		expiration = &expirationTime
	}

	if now.IsZero() {
		now = time.Now().UTC()
	}
	return service.cache.putKey(key, valueType, valueBytes, now, expiration)
}

func (service *Service) deleteKey(key string) error {
	now := time.Now().UTC()
	return service.cache.deleteKey(key, now)
}

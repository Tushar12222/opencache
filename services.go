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

func (service *Service) getKey(key string) (any, *time.Time, error) {
	now := time.Now().UTC()
	return service.cache.getKey(key, now)
}

func (service *Service) addKey(key string, value any, exp int) error {
	var expiration *time.Time
	if exp != 0 {
		now := time.Now().UTC()
		expirationTime := now.Add(time.Duration(exp) * time.Millisecond)
		expiration = &expirationTime
	}

	now := time.Now().UTC()
	return service.cache.addKey(key, value, now, expiration)
}

func (service *Service) putKey(key string, value any, exp int) error {
	var expiration *time.Time
	if exp != 0 {
		now := time.Now().UTC()
		expirationTime := now.Add(time.Duration(exp) * time.Millisecond)
		expiration = &expirationTime
	}

	now := time.Now().UTC()
	return service.cache.putKey(key, value, now, expiration)
}

func (service *Service) deleteKey(key string) error {
	now := time.Now().UTC()
	return service.cache.deleteKey(key, now)
}

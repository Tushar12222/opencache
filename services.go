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
	return service.cache.keyExists(key)
}

func (service *Service) getKey(key string) (any, *time.Time, error) {
	return service.cache.getKey(key)
}

func (service *Service) addKey(key string, value any, exp int) error {
	var expiration *time.Time
	if exp != 0 {
		now := time.Now().UTC()
		expirationTime := now.Add(time.Duration(exp) * time.Millisecond)
		expiration = &expirationTime
	}
	return service.cache.addKey(key, value, expiration)
}

func (service *Service) putKey(key string, value any, exp int) error {
	var expiration *time.Time
	if exp != 0 {
		now := time.Now().UTC()
		expirationTime := now.Add(time.Duration(exp) * time.Millisecond)
		expiration = &expirationTime
	}
	return service.cache.putKey(key, value, expiration)
}

func (service *Service) deleteKey(key string) error {
	return service.cache.deleteKey(key)
}

package main

import (
	"time"
)

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

func (service *Service) getKey(key string) (byte, []byte, *time.Time, *OCPError) {
	now := time.Now().UTC()
	valueType, valueBytes, expiry, err := service.cache.getKey(key, now)
	if err != nil {
		ocpError := parseErrorToOCPError(err)
		return 0, []byte{}, nil, ocpError
	}

	return valueType, valueBytes, expiry, nil
}

func (service *Service) addKey(key string, valueType byte, valueBytes []byte, exp int) *OCPError {
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

	err := service.cache.addKey(key, valueType, valueBytes, now, expiration)
	if err != nil {
		ocpError := parseErrorToOCPError(err)
		return ocpError
	}

	return nil
}

func (service *Service) putKey(key string, valueType byte, valueBytes []byte, exp int) *OCPError {
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

	err := service.cache.putKey(key, valueType, valueBytes, now, expiration)
	if err != nil {
		ocpError := parseErrorToOCPError(err)
		return ocpError
	}

	return nil
}

func (service *Service) deleteKey(key string) *OCPError {
	now := time.Now().UTC()

	err := service.cache.deleteKey(key, now)
	if err != nil {
		ocpError := parseErrorToOCPError(err)
		return ocpError
	}

	return nil
}

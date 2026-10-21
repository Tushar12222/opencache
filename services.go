package main

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

func (service *Service) getKey(key string) (any, error) {
	return service.cache.getKey(key)
}

func (service *Service) addKey(key string, value any, exp int) error {
	return service.cache.addKey(key, value)
}

func (service *Service) putKey(key string, value any, exp int) error {
	return service.cache.putKey(key, value)
}

func (service *Service) deleteKey(key string) error {
	return service.cache.deleteKey(key)
}

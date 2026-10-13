package main

type Service struct {
	cache *Cache
}

func (service *Service) getKey(key any) (any, error) {
	return service.cache.getKey(key)
}

func (service *Service) addKey(key any, value any) error {
	return service.cache.addKey(key, value)
}

func (service *Service) putKey(key any, value any) error {
	return service.cache.putKey(key, value)
}

func (service *Service) deleteKey(key any) error {
	return service.cache.deleteKey(key)
}

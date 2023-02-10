package cache

// Cache cache interface
type Cache interface {
	Get(key string) (any, error)
	Set(key string, val any)
}

var cache Cache = newRedis()

// Get get value in cache by key
func Get(key string) (any, error) {
	return cache.Get(key)
}

func Set(key string, val any) {
	cache.Set(key, val)
}

func FindStreamInfo(streamId string) (any, error) {
	return nil, nil
}

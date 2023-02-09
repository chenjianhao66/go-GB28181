package cache

// Cache cache interface
type Cache interface {
	Get(key string) (any, error)
	Set(key string, val any)
}

var cache Cache = newRedis()

func Get(key string) (any, error) {
	return cache.Get(key)
}

func Set(key string, val any) {
	cache.Set(key, val)
}

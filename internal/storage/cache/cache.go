package cache

// Cache cache interface
type Cache interface {
}

var cache Cache = newRedis()

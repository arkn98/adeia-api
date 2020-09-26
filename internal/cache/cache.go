package cache

// Cache is an interface for all cache-related functions, that implementations
// must implement.
type Cache interface {
	Close() error
	Get(dest interface{}, key string) error
	Set(key string, value string) error
}

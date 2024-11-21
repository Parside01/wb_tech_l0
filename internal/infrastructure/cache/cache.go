package cache

type Cacheable interface {
	Key() string
}

type Cache interface {
	Set(key string, value Cacheable)
	Get(key string) (Cacheable, bool)
	Delete(key string)
	Clear()
}

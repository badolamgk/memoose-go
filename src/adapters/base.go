package adapters

type CacheKey string
type TTL int64

type KeyValuePair[T any] struct {
	Key   CacheKey
	Value T
}

type CacheProviderCore[T any] interface {
	Get(key CacheKey) (*T, error)
	MGet(keys ...CacheKey) ([]*T, error)
	MSet(pairs ...KeyValuePair[T]) (string, error) // Assuming "OK" is returned as a string
	Set(key CacheKey, data T, ttl *TTL) (string, error)
	Del(keys ...CacheKey) (int, error)
	Expire(key CacheKey, newTTLFromNow TTL) (int, error)
}

type CacheProvider[T any] interface {
	CacheProviderCore[T]
	Name() string
	Pipeline() Pipeline[T]
	StoresAsObj() bool
}

type Pipeline[T any] interface {
	Get(key CacheKey) Pipeline[T]
	Set(key CacheKey, data T, ttl *TTL) Pipeline[T]
	Del(key CacheKey) Pipeline[T]
	Expire(key CacheKey, newTTLFromNow TTL) Pipeline[T]
	Exec() []any
}

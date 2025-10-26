package cache

import "time"

type RedisCacheService interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string, dest any) error
	Clear(key string) error
	Exists(key string) (bool, error)
}

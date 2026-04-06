package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type redisCacheService struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisCacheService(rdb *redis.Client) RedisCacheService {
	return &redisCacheService{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (r *redisCacheService) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.rdb.Set(r.ctx, key, data, ttl).Err()
}

func (r *redisCacheService) Get(key string, dest any) error {
	val, err := r.rdb.Get(r.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return err
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)

}

func (r *redisCacheService) Clear(key string) error {
	return r.rdb.Del(r.ctx, key).Err()
}

func (r *redisCacheService) Exists(key string) (bool, error) {
	count, err := r.rdb.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *redisCacheService) Incr(key string, ttl time.Duration) (int64, error) {
	val, err := r.rdb.Incr(r.ctx, key).Result()
	if err != nil {
		log.Printf("[Redis] Error Incr key %s: %v", key, err)
		return 0, err
	}

	if val == 1 {
		log.Printf("[Redis] New window started for key: %s (TTL: %v)", key, ttl)
		r.rdb.Expire(r.ctx, key, ttl)
	} else {
		log.Printf("[Redis] Incrementing key: %s, Current Count: %d", key, val)
	}

	return val, nil
}

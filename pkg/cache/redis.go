package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
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

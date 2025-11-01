package cache

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"
)

type RedisCache struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedisCache(addr string) *RedisCache {
    rdb := redis.NewClient(&redis.Options{
        Addr: addr, 
    })

    return &RedisCache{
        client: rdb,
        ctx:    context.Background(),
    }
}

func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return r.client.Set(r.ctx, key, data, ttl).Err()
}

func (r *RedisCache) Get(key string, dest interface{}) (bool, error) {
    val, err := r.client.Get(r.ctx, key).Result()
    if err == redis.Nil {
        return false, nil
    }
    if err != nil {
        return false, err
    }

    err = json.Unmarshal([]byte(val), dest)
    return true, err
}
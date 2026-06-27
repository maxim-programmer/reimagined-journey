package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const SearchTTL = 5 * time.Minute

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &RedisCache{client: client}
}

func (c *RedisCache) Ping(ctx context.Context) error {
	if err := c.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}
	return nil
}

func (c *RedisCache) Get(ctx context.Context, key string, dest any) (bool, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("redis get: %w", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, fmt.Errorf("redis unmarshal: %w", err)
	}
	return true, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis marshal: %w", err)
	}
	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}
	return nil
}

func (c *RedisCache) SetSession(ctx context.Context, token, userID string, ttl time.Duration) error {
	if err := c.client.Set(ctx, sessionKey(token), userID, ttl).Err(); err != nil {
		return fmt.Errorf("redis set session: %w", err)
	}
	return nil
}

func (c *RedisCache) GetSession(ctx context.Context, token string) (string, error) {
	val, err := c.client.Get(ctx, sessionKey(token)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("redis get session: %w", err)
	}
	return val, nil
}

func (c *RedisCache) DeleteSession(ctx context.Context, token string) error {
	if err := c.client.Del(ctx, sessionKey(token)).Err(); err != nil {
		return fmt.Errorf("redis delete session: %w", err)
	}
	return nil
}

func SearchKey(query string) string {
	return fmt.Sprintf("search:%s", query)
}

func sessionKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}
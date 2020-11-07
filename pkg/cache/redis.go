package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// Redis implements storage of objects in Redis, objects are serialised using JSON
type Redis struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedis returns a new instance of redis
func NewRedis(conn, pwd string, ttl time.Duration) (Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     conn,
		Password: pwd,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return Redis{}, errors.Wrap(err, "redis ping")
	}

	return Redis{
		client: client,
		ttl:    ttl,
	}, nil
}

// Get tries to retrieve an item from cache using the specified key. If the item does not exist, nil is returned.
func (c *Redis) Get(ctx context.Context, key string, typ interface{}) error {
	if typ == nil {
		return errors.New("return value `typ` must be set")
	}

	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return errors.Wrap(err, "redis get")
	}
	err = json.Unmarshal([]byte(result), &typ)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}

	return nil
}

// Set will add an item to the cache using the key. If an item exists, it is overwritten.
func (c *Redis) Set(ctx context.Context, key string, value interface{}) error {
	btsValue, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}
	err = c.client.Set(ctx, key, btsValue, c.ttl).Err()
	if err != nil {
		return errors.Wrap(err, "redis set")
	}

	return nil
}

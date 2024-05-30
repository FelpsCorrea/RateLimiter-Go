package limiter

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisClient{client: rdb}
}

func (r *RedisClient) Allow(ctx context.Context, key string, limit int, blockDuration time.Duration) (bool, error) {
	// Transaction to increment the count and set expiry
	txf := func(tx *redis.Tx) error {
		count, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		if count > limit {
			log.Printf("Limit reached for key: %s, count: %d, limit: %d", key, count, limit)
			return nil // Reached the limit
		}

		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Incr(ctx, key)
			pipe.Expire(ctx, key, blockDuration)
			return nil
		})

		log.Printf("Incremented count for key: %s, new count: %d", key, count+1)
		return err
	}

	// Run the transaction
	err := r.client.Watch(ctx, txf, key)
	if err != nil {
		return false, err
	}

	count, err := r.client.Get(ctx, key).Int()
	if err != nil {
		return false, err
	}

	log.Printf("Current count for key: %s, count: %d", key, count)
	return count <= limit, nil
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

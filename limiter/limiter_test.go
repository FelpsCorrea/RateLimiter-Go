package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupRedisTestClient() *RedisClient {
	client := NewRedisClient("redis:6379")
	return client
}

func TestRedisLimiter(t *testing.T) {
	lim := setupRedisTestClient()
	defer lim.Close()

	ctx := context.Background()

	key := "test:ip"
	limit := 5
	blockDuration := 10 * time.Second

	// Clear previous test data
	lim.client.Del(ctx, key)

	for i := 0; i < limit; i++ {
		allowed, err := lim.Allow(ctx, key, limit, blockDuration)
		assert.NoError(t, err)
		assert.True(t, allowed, "Request %d should be allowed", i+1)
	}

	// The next request should be blocked
	allowed, err := lim.Allow(ctx, key, limit, blockDuration)
	assert.NoError(t, err)
	assert.False(t, allowed, "Request should be blocked after reaching limit")

	// Wait for the block duration to expire and test again
	time.Sleep(blockDuration)

	allowed, err = lim.Allow(ctx, key, limit, blockDuration)
	assert.NoError(t, err)
	assert.True(t, allowed, "Request should be allowed after block duration expires")
}

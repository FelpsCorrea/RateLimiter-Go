package limiter

import (
	"context"
	"time"
)

type LimiterStorage interface {
	Allow(ctx context.Context, key string, limit int, blockDuration time.Duration) (bool, error)
	Close() error
}

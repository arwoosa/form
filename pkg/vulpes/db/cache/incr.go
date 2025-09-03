package cache

import (
	"context"
	"fmt"
)

func Incr(ctx context.Context, key string) (int64, error) {
	if conn == nil {
		return -1, ErrCacheNotConnected
	}
	val, err := conn.Incr(ctx, key).Result()
	if err != nil {
		return -1, fmt.Errorf("%w: %w", ErrCacheQueryFailed, err)
	}
	return val, nil
}

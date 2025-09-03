package cache

import (
	"context"
	"fmt"
)

func Keys(ctx context.Context, pattern string) ([]string, error) {
	if conn == nil {
		return nil, ErrCacheNotConnected
	}
	var err error
	var keys []string
	keys, err = conn.Keys(ctx, pattern).Result()
	if err == nil {
		return keys, nil
	}
	return nil, fmt.Errorf("%w: %w", ErrCacheQueryFailed, err)
}

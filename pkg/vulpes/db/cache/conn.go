package cache

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	conn *redis.Client
	once sync.Once

	defaultOptions = &redis.Options{
		PoolSize:     10,
		MinIdleConns: 3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
		IdleTimeout:  5 * time.Minute,
	}
)

type initConnOpt func(*redis.Options)

func WithAddr(addr string) initConnOpt {
	return func(o *redis.Options) {
		o.Addr = addr
	}
}

func WithDb(db int) initConnOpt {
	return func(o *redis.Options) {
		o.DB = db
	}
}

func WithPassword(password string) initConnOpt {
	return func(o *redis.Options) {
		o.Password = password
	}
}

func WithUsername(username string) initConnOpt {
	return func(o *redis.Options) {
		o.Username = username
	}
}

func InitConnection(opts ...initConnOpt) error {
	if conn != nil {
		return nil
	}
	once.Do(func() {
		for _, opt := range opts {
			opt(defaultOptions)
		}
		conn = redis.NewClient(defaultOptions)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if conn.Ping(ctx).Val() != "PONG" {
		return ErrCacheNotConnected
	}
	return nil
}

func Close() error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}

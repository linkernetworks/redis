package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

func NewPoolFromConfig(cf *RedisConfig) *redis.Pool {
	pool := redis.Pool{
		// the default max idle settings
		MaxIdle: 3,

		// the default idle timeout seconds
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cf.Addr())
		},
	}
	if cf.Pool != nil {
		if cf.Pool.MaxActive > 0 {
			pool.MaxActive = cf.Pool.MaxActive
		}
		if cf.Pool.MaxIdle > 0 {
			pool.MaxIdle = cf.Pool.MaxIdle
		}
		if cf.Pool.IdleTimeout > 0 {
			pool.IdleTimeout = cf.Pool.IdleTimeout * time.Second
		}
	}

	return &pool
}

// NewDefaultPool allocates the redis connection pool
func NewDefaultPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
}

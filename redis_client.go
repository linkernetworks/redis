package redis

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisService struct {
	Url  string
	Pool *redis.Pool
}

func NewPool(url string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", url) },
	}
}

func (service *RedisService) Do(cmd string, args ...interface{}) (interface{}, error) {
	c := service.Pool.Get()
	defer c.Close()
	return c.Do(cmd, args...)
}

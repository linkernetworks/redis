package redis

import (
	"bitbucket.org/linkernetworks/aurora/src/config"
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisService struct {
	Url  string
	Pool *redis.Pool
}

func (s *RedisService) Do(cmd string, args ...interface{}) (interface{}, error) {
	c := s.Pool.Get()
	defer c.Close()
	return c.Do(cmd, args...)
}

func NewService(cf *config.RedisConfig) *RedisService {
	addr := cf.Addr()
	return &RedisService{
		Url:  url,
		Pool: NewPool(addr),
	}
}

func NewPool(url string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", url) },
	}
}

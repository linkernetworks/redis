package redis

import (
	"bitbucket.org/linkernetworks/aurora/src/config"
	"github.com/garyburd/redigo/redis"
	"time"
)

type Service struct {
	Url  string
	Pool *redis.Pool
}

func (s *Service) Do(cmd string, args ...interface{}) (interface{}, error) {
	c := s.Pool.Get()
	defer c.Close()
	return c.Do(cmd, args...)
}

func NewService(cf *config.RedisConfig) *Service {
	addr := cf.Addr()
	return &Service{
		Url:  addr,
		Pool: NewPool(addr),
	}
}

func NewPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

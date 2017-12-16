package redis

import (
	"bitbucket.org/linkernetworks/aurora/src/config"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"time"
)

type Service struct {
	Url  string
	Pool *redis.Pool
}

func (s *Service) SetJSON(key string, m interface{}) error {
	c := s.Pool.Get()
	defer c.Close()

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err := c.Do("SET", key, bytes)
	return err
}

func (s *Service) PublishJSON(key string, m interface{}) error {
	c := s.Pool.Get()
	defer c.Close()

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err := c.Do("PUBLISH", key, bytes)
	return err
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

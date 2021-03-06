package redis

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

type Service struct {
	Pool *redis.Pool
}

// GetConnection returns the connection from the redis connection pool
func (s *Service) GetConnection() *Connection {
	conn := s.Pool.Get()
	return &Connection{conn}
}

func (s *Service) SetJSON(key string, m interface{}) error {
	c := s.GetConnection()
	defer c.Close()

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = c.Do("SET", key, bytes)
	return err
}

func (s *Service) PublishAndSetJSON(key string, m interface{}) error {
	if err := s.SetJSON(key, m); err != nil {
		return err
	}
	if err := s.PublishJSON(key, m); err != nil {
		return err
	}
	return nil
}

func (s *Service) Do(cmd string, args ...interface{}) (interface{}, error) {
	c := s.GetConnection()
	defer c.Close()
	return c.Do(cmd, args...)
}

func (s *Service) PublishJSON(key string, m interface{}) error {
	c := s.GetConnection()
	defer c.Close()

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = c.Do("PUBLISH", key, bytes)
	return err
}

// NewWithPool allocates a redis service with a given redis connection pool
func NewWithPool(pool *redis.Pool) *Service {
	return &Service{Pool: pool}
}

// New allocates a redis service with a given redis config
func New(cf *RedisConfig) *Service {
	if cf.Pool != nil {
		return &Service{Pool: NewPoolFromConfig(cf)}
	}
	return &Service{Pool: NewDefaultPool(cf.Addr())}
}

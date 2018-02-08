package redis

import (
	"encoding/json"
	"time"

	"bitbucket.org/linkernetworks/aurora/src/config"
	"github.com/garyburd/redigo/redis"
)

// TODO refactor: add logger
type Service struct {
	Pool *redis.Pool
}

// GetConnection returns the connection from the redis connection pool
func (s *Service) GetConnection() *Connection {
	conn := s.Pool.Get()
	return &Connection{conn}
}

// GetNumSub return a map of the subscriber with the redis channel key
// This method should be deprecated
func (s *Service) GetNumSub(key string) (map[string]int, error) {
	conn := s.GetConnection()
	defer conn.Close()
	return conn.PubSub().NumSub(key)
}

func (s *Service) SetJSON(key string, m interface{}) error {
	c := s.Pool.Get()
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

func (s *Service) PublishJSON(key string, m interface{}) error {
	c := s.Pool.Get()
	defer c.Close()

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = c.Do("PUBLISH", key, bytes)
	return err
}

func (s *Service) Do(cmd string, args ...interface{}) (interface{}, error) {
	c := s.Pool.Get()
	defer c.Close()
	return c.Do(cmd, args...)
}

func New(cf *config.RedisConfig) *Service {
	return &Service{Pool: NewPoolFromConfig(cf)}
}

func NewPoolFromConfig(cf *config.RedisConfig) *redis.Pool {
	pool := redis.Pool{
		// the default max idle settings
		MaxIdle: 3,

		// the default idle timeout seconds
		IdleTimeout: 240 * time.Second,

		MaxActive: 500,

		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cf.Addr())
		},
	}
	if cf.Pool.MaxActive > 0 {
		pool.MaxActive = cf.Pool.MaxActive
	}
	if cf.Pool.MaxIdle > 0 {
		pool.MaxIdle = cf.Pool.MaxIdle
	}
	if cf.Pool.IdleTimeout > 0 {
		pool.IdleTimeout = cf.Pool.IdleTimeout * time.Second
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

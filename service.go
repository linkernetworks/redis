package redis

import (
	"encoding/json"

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
	// return &Service{Pool: NewPoolFromConfig(cf)}
	return &Service{Pool: NewDefaultPool(cf.Addr())}
}

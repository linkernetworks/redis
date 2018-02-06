package redis

import (
	"encoding/json"
	"time"

	"bitbucket.org/linkernetworks/aurora/src/config"
	"github.com/garyburd/redigo/redis"
)

// TODO refactor: add logger
type Service struct {
	Url  string
	Pool *redis.Pool
}

// GetNumSub return a map of the subscriber with the redis channel key
func (s *Service) GetNumSub(key string) (map[string]int, error) {
	c := s.Pool.Get()
	defer c.Close()
	return redis.IntMap(c.Do("PUBSUB", "NUMSUB", key))
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

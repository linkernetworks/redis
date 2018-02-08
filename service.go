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

type Connection struct {
	redis.Conn
}

// GetConnection returns the connection from the redis connection pool
func (s *Service) GetConnection() *Connection {
	conn := s.Pool.Get()
	return &Connection{conn}
}

// GetString executes the command "GET" with a given key, and cast the result
// to go string type
func (c *Connection) GetString(key string) (val string, err error) {
	val, err = redis.String(c.Do("GET", key))
	return val, err
}

// GetInt executes the command "GET" with a given key, and cast the result
// to go int type
func (c *Connection) GetInt(key string) (val int, err error) {
	val, err = redis.Int(c.Do("GET", key))
	return val, err
}

// SetJSON executes the command "SET" and encode the json into bytes to set
// with the given key.
func (c *Connection) SetJSON(key string, m interface{}) (err error) {
	var bytes []byte
	bytes, err = json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = c.Do("SET", key, bytes)
	return err
}

func (c *Connection) PublishAndSetJSON(key string, m interface{}) (err error) {
	if err = c.SetJSON(key, m); err != nil {
		return err
	}
	if err = c.PublishJSON(key, m); err != nil {
		return err
	}
	return err
}

func (c *Connection) PublishJSON(key string, m interface{}) (err error) {
	var bytes []byte
	bytes, err = json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = c.Do("PUBLISH", key, bytes)
	return err
}

func (c *Connection) PubSub() *PubSubCommandContext {
	return &PubSubCommandContext{c, "PUBSUB"}
}

type PubSubCommandContext struct {
	Connection *Connection
	Command    string
}

func (cc *PubSubCommandContext) NumSub(key string) (m map[string]int, err error) {
	m, err = redis.IntMap(cc.Connection.Do(cc.Command, "NUMSUB", key))
	return m, err
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

package redis

import (
	"encoding/json"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Connection struct {
	redis.Conn
}

func (c *Connection) SetWithExpire(key string, value interface{}, expire int) (reply interface{}, err error) {
	return c.Do("SET", key, value, "EX", expire)
}

func (c *Connection) Set(key string, value interface{}) (reply interface{}, err error) {
	return c.Do("SET", key, value)
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

func (c *Connection) Ping() (reply interface{}, err error) {
	reply, err = c.Do("PING")
	return reply, err
}

// KeepAlive uses a ticker to ping the connection, this function returns a
// time.Ticker object.  Be sure to stop the ticker when the connection is not
// being used.
func (c *Connection) KeepAlive(interval time.Duration) *time.Ticker {
	var ticker = time.NewTicker(interval)
	go func() {
		for t := range ticker.C {
			if _, err := c.Ping(); err != nil {
				log.Printf("redis: failed to ping error=%v at %s\n", err, t.String())
			}
		}
	}()
	return ticker
}

func (c *Connection) PubSub() *PubSubCommandContext {
	return &PubSubCommandContext{c, &redis.PubSubConn{c}}
}

func (c *Connection) ZSet(key string) *ZSet {
	return NewZSet(c, key)
}

type PubSubCommandContext struct {
	*Connection
	*redis.PubSubConn
}

func (cc *PubSubCommandContext) NumSub(key string) (m map[string]int, err error) {
	m, err = redis.IntMap(cc.Do("PUBSUB", "NUMSUB", key))
	return m, err
}

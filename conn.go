package redis

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

type Connection struct {
	redis.Conn
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

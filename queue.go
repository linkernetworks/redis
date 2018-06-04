package redis

import (
	"encoding/json"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
)

// Queue: Task queue to process with key
type Queue struct {
	*Connection
}

//New task queue
func NewQueue(conn *Connection) *Queue {
	return &Queue{
		Connection: conn,
	}
}

// EnqueueString with key
func (q *Queue) EnqueueString(key string, value string) (n int, err error) {
	return redigo.Int(q.Do("RPUSH", key, value))
}

// EnqueueJSON with key, value could be interface object with JSON
func (q *Queue) EnqueueJSON(key string, value interface{}) (n int, err error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return 0, err
	}

	return redigo.Int(q.Do("RPUSH", key, bytes))
}

// DequeueString with key
func (q *Queue) DequeueString(key string) (string, error) {
	return redigo.String(q.Do("LPOP", key))
}

// Dequeue with key
func (q *Queue) DequeueJSON(key string) ([]byte, error) {
	val, err := q.Do("LPOP", key)
	if err != nil {
		return nil, err
	}
	switch v := val.(type) {
	case []byte:
		return v, nil
	}

	return nil, fmt.Errorf("%s", "Error in type assertion")
}

// Cleanup all data in specific key in queue
func (q *Queue) RemoveAll(key string) (int, error) {
	return redigo.Int(q.Do("DEL", key, -1, 0))
}

func (q *Queue) Len(key string) (int, error) {
	return redigo.Int(q.Do("LLEN", key))
}

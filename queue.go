package redis

import redigo "github.com/gomodule/redigo/redis"

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

// Enqueue with key
func (q *Queue) Enqueue(key string, value interface{}) (n int, err error) {
	return redigo.Int(q.Do("RPUSH", key, value))
}

// Dequeue with key
func (q *Queue) Dequeue(key string) (string, error) {
	return redigo.String(q.Do("LPOP", key))
}

// Cleanup all data in specific key in queue
func (q *Queue) RemoveAll(key string) (int, error) {
	return redigo.Int(q.Do("DEL", key, -1, 0))
}

func (q *Queue) Len(key string) (int, error) {
	return redigo.Int(q.Do("LLEN", key))
}

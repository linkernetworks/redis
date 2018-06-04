package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleQueue(t *testing.T) {
	redisConfig := &RedisConfig{
		Host: "localhost",
		Port: 6379,
	}

	redisService := New(redisConfig)

	conn := redisService.GetConnection()
	defer conn.Close()

	q := NewQueue(conn)
	assert.NotNil(t, q)

	_, err := q.RemoveAll("test1")
	assert.NoError(t, err)

	n, err := q.Enqueue("test1", "value1")
	assert.NoError(t, err)
	t.Log(n)
	assert.NotZero(t, n)

	n2, err2 := q.Enqueue("test1", "value2")
	assert.NoError(t, err2)
	t.Log(n2)
	assert.NotZero(t, n2)

	v1, err3 := q.Dequeue("test1")
	assert.NoError(t, err3)
	assert.Equal(t, "value1", v1)

	v2, err4 := q.Dequeue("test1")
	assert.NoError(t, err4)
	assert.Equal(t, "value2", v2)
}

func TestMultipleQueue(t *testing.T) {
	redisConfig := &RedisConfig{
		Host: "localhost",
		Port: 6379,
	}

	redisService := New(redisConfig)

	conn := redisService.GetConnection()
	defer conn.Close()

	q := NewQueue(conn)
	assert.NotNil(t, q)

	_, err := q.RemoveAll("test1")
	assert.NoError(t, err)

	_, err = q.RemoveAll("test2")
	assert.NoError(t, err)

	n, err := q.Enqueue("test1", "value1")
	assert.NoError(t, err)
	assert.NotZero(t, n)

	n, err = q.Enqueue("test2", "value3")
	assert.NoError(t, err)
	assert.NotZero(t, n)

	n, err = q.Enqueue("test1", "value2")
	assert.NoError(t, err)
	assert.NotZero(t, n)

	n, err = q.Enqueue("test2", "value4")
	assert.NoError(t, err)
	assert.NotZero(t, n)

	v1, err := q.Dequeue("test1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", v1)

	v2, err := q.Dequeue("test1")
	assert.NoError(t, err)
	assert.Equal(t, "value2", v2)

	v3, err := q.Dequeue("test2")
	assert.NoError(t, err)
	assert.Equal(t, "value3", v3)

	v4, err := q.Dequeue("test2")
	assert.NoError(t, err)
	assert.Equal(t, "value4", v4)
}

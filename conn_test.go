package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetKeyValue(t *testing.T) {
	redisConfig := &RedisConfig{
		Host: "localhost",
		Port: 6379,
	}

	redisService := New(redisConfig)

	conn := redisService.GetConnection()
	defer conn.Close()

	_, err := conn.Set("key1", "value")
	assert.NoError(t, err)

	v, err := conn.GetString("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value", v)

	_, err = conn.Set("key2", 42)
	assert.NoError(t, err)

	v2, err := conn.GetInt("key2")
	assert.NoError(t, err)
	assert.Equal(t, 42, v2)
}

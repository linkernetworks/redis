package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testSetKeyValue(t *testing.T) {
	redisConfig := &RedisConfig{
		Host: "localhost",
		Port: 6379,
	}

	redisService := New(redisConfig)

	conn := redisService.GetConnection()
	defer conn.Close()

	_, err := conn.Set("key", "value")
	assert.NoError(t, err)

	v, err := conn.GetString("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", v)

}

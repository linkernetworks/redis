package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisConfig(t *testing.T) {

	inf := "testinterface"
	cf := RedisConfig{
		Interface: inf,
	}
	assert.True(t, cf.Unresolved())

	err := cf.LoadDefaults()
	assert.Equal(t, cf.Port, int32(6379))
	assert.Equal(t, cf.Host, "localhost")
	assert.NoError(t, err)

	host := "imhost"
	port := 6380
	addr := "imhost:6380"
	cf.SetHost(host)
	cf.SetPort(int32(port))
	assert.Equal(t, cf.Addr(), addr)
	assert.Nil(t, cf.GetPublic())
	assert.Equal(t, cf.GetInterface(), inf)
}

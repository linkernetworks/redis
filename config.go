package redis

import (
	"bitbucket.org/linkernetworks/aurora/src/config/serviceconfig"
	"net"
	"strconv"
	"time"
)

type RedisPoolConfig struct {
	// max number of the idle connections
	MaxIdle int `json:"maxIdle"`

	// max number of the active connections
	MaxActive int `json:"maxActive"`

	// the timeout seconds for idle connections
	IdleTimeout time.Duration `json:"idleTimeout"`
}

type RedisConfig struct {
	Host string `json:"host"`
	Port int32  `json:"port"`

	// net interface for dynamically assign the IP
	Interface string `json:"interface"`

	Pool *RedisPoolConfig `json:"pool"`

	Public *RedisConfig `json:"public"`
}

func (c *RedisConfig) Unresolved() bool {
	return c.Host == ""
}

func (c *RedisConfig) SetHost(host string) {
	c.Host = host
}

func (c *RedisConfig) SetPort(port int32) {
	c.Port = port
}

// Implement DefaultLoader
func (c *RedisConfig) LoadDefaults() error {
	if c.Port == 0 {
		c.Port = 6379
	}
	if c.Host == "" {
		c.Host = "localhost"
	}
	return nil
}

func (c *RedisConfig) GetInterface() string {
	return c.Interface
}

func (c *RedisConfig) GetPublic() serviceconfig.ServiceConfig {
	return c.Public
}

// Addr implements the Address interface
func (c *RedisConfig) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(int(c.Port)))
}

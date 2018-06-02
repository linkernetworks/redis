Redis
===

[![Build Status](https://travis-ci.org/linkernetworks/redis.svg?branch=master)](https://travis-ci.org/linkernetworks/redis)

[![codecov](https://codecov.io/gh/linkernetworks/redis/branch/master/graph/badge.svg)](https://codecov.io/gh/linkernetworks/redis)

Redis is a package integrating redisDB with redigo.

# How to use

##### Example

```
redisConfig := &redis.RedisConfig{
  Host: "localhost",
  Port: 6379, 
}

redisService := redis.New(redisConfig)

conn := redisService.GetConnection()
defer conn.Close()

if _, err := conn.Set("key", "value"); err != nil {
  // handle err
}

v, err := conn.GetString("key")
 
```

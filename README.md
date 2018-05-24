Redis
===

[![Build Status](https://travis-ci.org/linkernetworks/redis.svg?branch=master)](https://travis-ci.org/linkernetworks/redis)

Redis is a package integrating redisDB with redigo.

# How to use

##### Example

```
const redisConfig := *redis.RedisConfig{
  Host: "localhost",
  Port: 6379, 
}

redisService := redis.New(redisConfig)

conn := redis.GetConnection()
defer conn.Close()

if err := conn.Set("key", "value"); err != nil {
  // handle err
}

v, err := conn.GetString("key")
 
```

package redis

import (
	"errors"

	redigo "github.com/garyburd/redigo/redis"

	"bitbucket.org/linkernetworks/aurora/src/jobcontroller/types"
)

// constant error messages are used for i18n
const (
	ErrConvertRedisResponseMsg = "convert redis response"
)

var (
	// ErrConvertRedisResponse returned when convert Redis response error (type assertion)
	ErrConvertRedisResponse = errors.New(ErrConvertRedisResponseMsg)
)

// ZSet is a client of Redis sorted set (ZSET for short).
// https://redis.io/commands#sorted_set
type ZSet struct {
	rds *Connection
	// Key is name of ZSET
	key string
	// Mtx is a read/write lock for Redis
	// mtx *sync.RWMutex
}

// NewZSet creates a RedisZSet with internal fields initialized
func NewZSet(conn *Connection, zsetKey string) *ZSet {
	return &ZSet{
		rds: conn,
		key: zsetKey,
	}
}

// ZADD add an KV(score as key, member as value) to Redis.
// Return n means number of elements added to the sorted sets.
// Return err if any error occured.
// See https://redis.io/commands/zadd
func (rz *ZSet) ZADD(score float64, member interface{}) (n int, err error) {
	// ZADD key [NX|XX] [CH] [INCR] score member [score member ...]
	return redigo.Int(rz.rds.Do("ZADD", rz.key, score, member))
}

// ZRANGEBYSCORE ranges over ZSET ( where  min < score && score < max )
// See https://redis.io/commands/zrangebyscore
func (rz *ZSet) ZRANGEBYSCORE(min, max float64, offset, limit int) (members []interface{}, err error) {
	// ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]
	return redigo.Values(rz.rds.Do("ZRANGEBYSCORE", rz.key, min, max, "LIMIT", offset, limit))
}

// Len returns length of ZSET elements
func (rz *ZSet) Len() int {
	len, err := redigo.Int(rz.rds.Do("ZCARD", rz.key))
	if err != nil {
		return 0 // nothing in the ZSET or key not exist
	}
	return len
}

// ZREM removes one member from Redis ZSET
// Return n means number of elements removed from the sorted sets.
// Return err if any error occured.
// See https://redis.io/commands/zrem
func (rz *ZSet) ZREM(member interface{}) (int, error) {
	return redigo.Int(rz.rds.Do("ZREM", rz.key, member))
}

// RemoveAll drops all data in a Redis ZSET, use with CAUTION
// See https://redis.io/commands/zremrangebyscore
func (rz *ZSet) RemoveAll() (int, error) {
	return redigo.Int(rz.rds.Do("ZREMRANGEBYSCORE", rz.key, "-inf", "+inf"))
}

// ZPOP pops a value from the ZSET key using ZRANGEBYSCORE/ZREM commands.
// TODO sort by enque time
// TODO benchmark
// TODO need transaction?
func (rz *ZSet) ZPOP() (interface{}, error) {
	members, err := rz.ZRANGEBYSCORE(types.ScoreHigh, types.ScoreLow, 0, 1)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, redigo.ErrNil
	}
	if _, err = rz.ZREM(members[0]); err != nil {
		return nil, err
	}
	return members[0], nil
}

func (rz *ZSet) QueryAll() (members []interface{}, err error) {
	return redigo.Values(rz.rds.Do("ZRANGEBYSCORE", rz.key, "-inf", "+inf"))
}
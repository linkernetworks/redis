package redis

type RedisConfig struct {
	Host      string       `json:"host"`
	Port      int          `json:"port"`
	Interface string       `json:"interface"`
	Public    *RedisConfig `json:"public"`
}

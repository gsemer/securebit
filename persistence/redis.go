package persistence

import (
	"fmt"
	"securebit/config"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host string
	Port string
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host: config.GetEnv("REDIS_HOST", ""),
		Port: config.GetEnv("REDIS_PORT", ""),
	}
}

func (config *RedisConfig) Init() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.Host, config.Port),
		Password: "",
		DB:       0,
	})
}

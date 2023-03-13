package cache

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

type IRedisClient interface {
	Get(key string) RedisStringCmd
	Set(key string, value interface{}, expiration time.Duration) StatusCmd
}

type RedisStringCmd interface {
	Uint64() (uint64, error)
	Err() error
}

type StatusCmd interface {
	Err() error
}

type redisClientImpl struct {
	client *redis.Client
}

func NewRedisClient() IRedisClient {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	if redisClient == nil {
		log.Panicln("can not connect redis")
	}
	return &redisClientImpl{client: redisClient}
}

func (r *redisClientImpl) Get(key string) RedisStringCmd {
	cmd := r.client.Get(key)
	return &RedisStringCmdImpl{cmd}
}

func (r *redisClientImpl) Set(key string, value interface{}, expiration time.Duration) StatusCmd {
	cmd := r.client.Set(key, value, expiration)
	return &StatusCmdImpl{cmd}
}

type RedisStringCmdImpl struct {
	cmd *redis.StringCmd
}

func (r *RedisStringCmdImpl) Uint64() (uint64, error) {
	return r.cmd.Uint64()
}

func (r *RedisStringCmdImpl) Err() error {
	return r.cmd.Err()
}

type StatusCmdImpl struct {
	cmd *redis.StatusCmd
}

func (s *StatusCmdImpl) Err() error {
	return s.cmd.Err()
}

// ---------

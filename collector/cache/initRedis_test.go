package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	//cache "collector/cache"
	mock "collector/mocks"

	"github.com/go-redis/redis"
)

//func TestNewRedisClient(t *testing.T) {
//	// Test successful creation of Redis client
//	client := cache.NewRedisClient()
//	assert.NotNil(t, client)
//
//	// Test panic if Redis client can't connect
//	redis.NewClient = func(options *redis.Options) *redis.Client {
//		return nil
//	}
//
//	defer func() {
//		redis.NewClient = redis.NewClient
//	}()
//	assert.Panics(t, func() {
//		cache.NewRedisClient()
//	})
//}

func TestRedisClientImpl_Get(t *testing.T) {
	mockRedisStringCmd := &mock.RedisStringCmd{}
	mockRedisClient := &mock.IRedisClient{}
	key := "test-key"

	// Test successful Get operation
	mockRedisClient.On("Get", key).Return(mockRedisStringCmd)
	client := &redisClientImpl{client: mockRedisClient}
	cmd := client.Get(key)
	assert.NotNil(t, cmd)
	assert.Equal(t, mockRedisStringCmd, cmd.(*cache.RedisStringCmdImpl).Cmd)

	// Test Get operation that returns an error
	mockRedisClient.On("Get", key).Return(nil, redis.Nil)
	cmd = client.Get(key)
	assert.NotNil(t, cmd)
	assert.Equal(t, redis.Nil, cmd.(*cache.RedisStringCmdImpl).Err())
}

func TestRedisClientImpl_Set(t *testing.T) {
	mockStatusCmd := &redis.StatusCmd{}
	mockRedisClient := &mock.RedisClient{}
	key := "test-key"
	value := "test-value"
	expiration := time.Minute

	// Test successful Set operation
	mockRedisClient.On("Set", key, value, expiration).Return(mockStatusCmd)
	client := &cache.RedisClientImpl{Client: mockRedisClient}
	cmd := client.Set(key, value, expiration)
	assert.NotNil(t, cmd)
	assert.Equal(t, mockStatusCmd, cmd.(*cache.StatusCmdImpl).Cmd)

	// Test Set operation that returns an error
	mockRedisClient.On("Set", key, value, expiration).Return(nil, redis.Nil)
	cmd = client.Set(key, value, expiration)
	assert.NotNil(t, cmd)
	assert.Equal(t, redis.Nil, cmd.(*cache.StatusCmdImpl).Err())
}

package cache

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockRedisClient struct {
	mock.Mock
	testing *testing.T
}

func (m *MockRedisClient) Get(key string) RedisStringCmd {
	args := m.Called(key)
	return args.Get(0).(RedisStringCmd)
}

func (m *MockRedisClient) Set(key string, value interface{}, expiration time.Duration) StatusCmd {
	args := m.Called(key, value, expiration)
	return args.Get(0).(StatusCmd)
}

func TestNewRedisClient(t *testing.T) {
	mockRedisClient := new(MockRedisClient)

	// Thiết lập mong muốn của mock object
	mockRedisClient.On("Get", "test").Return(new(RedisStringCmdImpl))
	mockRedisClient.On("Set", "test", "test", 1*time.Second).Return(new(StatusCmdImpl))

	// Gọi hàm để được test
	redisClient := &redisClientImpl{
		client: mockRedisClient,
	}
	assert.NotNil(t, redisClient)

	// Kiểm tra xem mock object đã được gọi đúng chưa
	mockRedisClient.AssertExpectations(t)
}

func TestRedisClientImpl_Get(t *testing.T) {
	mockRedisClient := new(MockRedisClient)

	// Thiết lập mong muốn của mock object
	mockRedisClient.On("Get", "test").Return(new(RedisStringCmdImpl))

	redisClient := &redisClientImpl{
		client: mockRedisClient,
	}
	assert.NotNil(t, redisClient)

	// Gọi phương thức để được test
	cmd := redisClient.Get("test")
	assert.NotNil(t, cmd)

	// Kiểm tra xem mock object đã được gọi đúng chưa
	mockRedisClient.AssertExpectations(t)
}

func TestRedisClientImpl_Set(t *testing.T) {
	mockRedisClient := new(MockRedisClient)
	// Thiết lập mong muốn của mock object
	mockRedisClient.On("Set", "test", "test", 1*time.Second).Return(new(StatusCmdImpl))

	redisClient := &redisClientImpl{
		client: mockRedisClient,
	}
	assert.NotNil(t, redisClient)

	// Gọi phương thức để được test
	cmd := redisClient.Set("test", "test", 1*time.Second)
	assert.NotNil(t, cmd)

	// Kiểm tra xem mock object đã được gọi đúng chưa
	mockRedisClient.AssertExpectations(t)
}

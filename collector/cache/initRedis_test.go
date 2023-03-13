package cache

import (
	"collector/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewRedisClient(t *testing.T) {
	mockClient := new(mocks.IRedisClient)
	redisClient := mocks.NewIRedisClient(mockClient)

	assert.NotNil(t, redisClient)

	mockClient.On("Get", "test").Return(&mocks.RedisStringCmd{}, nil)
	mockClient.On("Set", "test", "value", time.Duration(0)).Return(&mocks.StatusCmd{}, nil)

	_, err := redisClient.Get("test").Uint64()
	assert.Nil(t, err)

	err = redisClient.Set("test", "value", time.Duration(0)).Err()
	assert.Nil(t, err)

	mockClient.AssertExpectations(t)
}

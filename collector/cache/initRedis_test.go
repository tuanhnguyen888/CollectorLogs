package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewRedisClient(t *testing.T) {
	redisClient := NewRedisClient()
	if r, ok := redisClient.(IRedisClient); !ok {
		t.Errorf("Type interface IRedisClient but not %T", r)
	}
}

func TestRedisClientImpl_GetSet(t *testing.T) {
	redisClient := NewRedisClient()
	//1
	err := redisClient.Set("test", uint64(123), 10*time.Hour).Err()
	assert.Nil(t, err)
	//2
	number, err := redisClient.Get("test").Uint64()
	assert.Nil(t, err)
	if number != uint64(123) {
		t.Errorf("Expect  uint64 123 but is %T %v", number, number)
	}
	//3
	err = redisClient.Get("test").Err()
	assert.Nil(t, err)
	//4
	err = redisClient.Get("abcdefgh").Err()
	assert.NotNil(t, err)
}

package database

import (
	"github.com/go-redis/redis"
	"testing"
)

func TestNewCassandra(t *testing.T) {
	// Tạo một mock object cho redis.Client
	redisMock := &redis.Client{}

	// Gọi hàm NewCassandra
	db := NewCassandra("127.0.0.1", "query", redisMock)

	// Kiểm tra xem đối tượng redisConn của cassandraStory có được gán bằng đối tượng redisMock không
	if db.(*cassandraStory).redisConn != redisMock {
		t.Errorf("redis connection not set correctly")
	}
}

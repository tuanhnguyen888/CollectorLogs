package database_test

import (
	"collector/database"
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestNewCassandra(t *testing.T) {
//	// Tạo một mock object cho redis.Client
//	redisMock := &redis.Client{}
//
//	// Gọi hàm NewCassandra
//	db := NewCassandra("127.0.0.1", "query", redisMock)
//
//	// Kiểm tra xem đối tượng redisConn của cassandraStory có được gán bằng đối tượng redisMock không
//	if db.(*cassandraStory).redisConn != redisMock {
//		t.Errorf("redis connection not set correctly")
//	}
//}

func TestNewCassandraConn(t *testing.T) {
	dns := "127.0.0.1:9042"
	db, err := database.NewCassandraConn(dns)
	assert.Nil(t, err)

	results := make(map[string]interface{})

	err = db.Query("SELECT * FROM collector.logs LIMIT 1 ALLOW FILTERING").Consistency(gocql.One).MapScan(results)
	assert.Nil(t, err)

	//var logs []map[string]interface{}
	_, err = db.Query("SELECT * FROM collector.logs LIMIT 1 ALLOW FILTERING").Iter().SliceMap()
	assert.Nil(t, err)

}

package database

import (
	redis2 "github.com/go-redis/redis"
	"testing"
)

func TestNewMssql(t *testing.T) {
	dsn := "sqlserver://SA:Khong123@127.0.0.1?database=TestDB"
	query := "SELECT * FROM logs"
	var redis *redis2.Client
	db := NewMssql(dsn, query, redis)

	if db == nil {
		t.Errorf("Expected NewPostgres() to return a non-nil value, but got nil")
	}

	// Kiểm tra kiểu trả về
	if _, ok := db.(*mssqlStory); !ok {
		t.Errorf("Expected NewPostgres() to return a *postgresqlStory, but got %T", db)
	}

	// Kiểm tra giá trị các thuộc tính
	ps, _ := db.(*mssqlStory)

	if ps.query != query {
		t.Errorf("Expected query = %q, but got %q", query, ps.query)
	}
}

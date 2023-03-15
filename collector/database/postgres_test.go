package database_test

import (
	"collector/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGormPgDB(t *testing.T) {
	dns := "host=127.0.0.1 port=5432 user=postgres password=khong123 dbname=user sslmode=disable"
	db := database.NewGormPgDB(dns)
	if _, ok := db.(database.IDb); !ok {
		t.Errorf("Exepect Interface IDb but not %T", db)
	}

	var logs []map[string]interface{}
	err := db.Raw("SELECT * FROM logs limit 1").Scan(&logs)
	assert.Nil(t, err)

	row, err := db.Raw("SELECT * FROM logs limit 1").Rows()
	assert.Nil(t, err)

	_, err = row.Columns()
	assert.Nil(t, err)

	err = row.Close()
	assert.Nil(t, err)

	boo := row.Next()
	if boo {
		t.Errorf("Expect True but not false")
	}
}

//func TestExecuteErrorRedisError(t *testing.T) {
//	mockDB := new(mocks.IDb)
//
//	mockRedis := new(mocks.IRedisClient)
//
//	mockRedis.On("Get", "timestamp_1").Return(nil, errors.New("error in Redis"))
//
//	db := database.NewPostgres("SELECT * FROM projects", mockDB, mockRedis)
//
//	_, err := db.Execute(1)
//	if assert.Error(t, err) {
//		assert.Equal(t, "error in Redis", err.Error())
//	}
//}

//func TestPostgresqlStory_Execute(t *testing.T) {
//	// Create mock DB
//	mockDB := new(mocks.IDb)
//	mockIRaw := new(mocks.IRaw)
//	mockIRow := new(mocks.IRow)
//
//	// Create mock Redis client
//	mockRedisClient := new(mocks.IRedisClient)
//
//	// Create test data
//	query := "SELECT * FROM test_column"
//	stt := 123
//	nameTable := "test_table"
//	columnName := "test_column"
//	timestamp := uint64(1234567890)
//	logs := []map[string]interface{}{{"id": 1, "name": "test1", "time": 1234567890}, {"id": 2, "name": "test2", "time": 1234567990}}
//	logsJSON, _ := json.Marshal(logs)
//
//	// Set up expectations for mock DB
//	mockDB.On("Raw", query+" LIMIT 1").Return(mockIRaw)
//
//	mockIRaw.On("Rows").Return(mockIRow, nil)
//	mockIRow.On("Close").Return(nil)
//	mockIRow.On("Columns").Return([]string{"id", "name", "time"}, nil)
//	mockIRow.On("Next").Return(true).Times(1)
//	mockIRow.On("Scan").Return(nil)
//
//	mockDB.On("Raw", "SELECT MAX("+columnName+") FROM "+nameTable).Return(mockIRow)
//
//	mockDB.On("Scan", &timestamp).Return(nil)
//	mockDB.On("Raw", "SELECT * FROM "+nameTable+" where "+columnName+" > ?", timestamp).Return(mockDB)
//	mockDB.On("Scan", mock.AnythingOfType("*[]map[string]interface {}")).Return(nil)
//
//	// Set up expectations for mock Redis client
//	mockRedisClient.On("Get", "timestamp_"+strconv.Itoa(stt)).Return(redis.NewStringResult(strconv.Itoa(int(timestamp)), nil))
//	mockRedisClient.On("Set", "timestamp_"+strconv.Itoa(stt), mock.Anything, 10*time.Hour).Return(nil)
//
//	// Create PostgresqlStory object with mock DB and Redis client
//	ps := database.NewPostgres(query, mockDB, mockRedisClient)
//
//	// Call the function being tested
//	result, err := ps.Execute(stt)
//
//	// Check the result
//	assert.Equal(t, logsJSON, result)
//	assert.Nil(t, err)
//
//	// Assert that the expectations for the mock objects were met
//	//mockDB.AssertExpectations(t)
//	//mockRedisClient.AssertExpectations(t)
//}

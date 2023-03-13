package database

import (
	"collector/cache"
	"collector/kafka"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

type ICassandraDatabase interface {
	Query(query string, args ...interface{}) IQuery
}

type IQuery interface {
	Consistency(c gocql.Consistency) IQuery
	MapScan(m map[string]interface{}) error
	Iter() IIter
	Scan(dest *uint64) error
}

type IIter interface {
	SliceMap() ([]map[string]interface{}, error)
}

type query struct {
	query *gocql.Query
}

type iter struct {
	iter *gocql.Iter
}

type CassandraConn struct {
	session *gocql.Session
}

func (c *CassandraConn) Query(queryStr string, args ...interface{}) IQuery {

	queryGocql := c.session.Query(queryStr, args...)
	query := new(query)
	query.query = queryGocql
	return query
}

func (q *query) Consistency(c gocql.Consistency) IQuery {
	queryGocql := q.query.Consistency(gocql.One)
	q.query = queryGocql
	return q
}

func (q *query) MapScan(m map[string]interface{}) error {
	return q.query.MapScan(m)
}
func (q *query) Iter() IIter {
	iterGocql := q.query.Iter()
	iterSt := new(iter)
	iterSt.iter = iterGocql
	return iterSt
}
func (q *query) Scan(dest *uint64) error {
	return q.query.Scan(&dest)
}

func (i *iter) SliceMap() ([]map[string]interface{}, error) {
	return i.iter.SliceMap()
}

func NewCassandraConn(dns string) (ICassandraDatabase, error) {
	cluster := gocql.NewCluster(dns)
	cluster.Keyspace = "collector"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &CassandraConn{session: session}, nil
}

type cassandraStory struct {
	query     string
	db        ICassandraDatabase
	redisConn cache.IRedisClient
}

func NewCassandra(query string, cassandra ICassandraDatabase, conn cache.IRedisClient) Database {
	return &cassandraStory{
		query:     query,
		db:        cassandra,
		redisConn: conn,
	}
}

func (s *cassandraStory) Execute(stt int) ([]byte, error) {
	var logs []map[string]interface{}
	sttString := strconv.Itoa(stt)

	//dns := s.redisClient.Get("dns_" + sttString).String()
	// Lấy một bản ghi từ bảng users
	results := make(map[string]interface{})
	if err := s.db.Query(s.query + " LIMIT 1 AlLOW FILTERING").Consistency(gocql.One).MapScan(results); err != nil {
		panic(err)
	}
	var columnName string
	// In ra thông tin về bản ghi
	for key, value := range results {
		valueInt, ok := value.(int64)
		if !ok {
			continue
		}
		if len(strconv.FormatInt(valueInt, 10)) == 13 || len(strconv.FormatInt(valueInt, 10)) == 11 {
			columnName = key
			break
		}
	}
	//log.Println(columnName, " casandra")
	timestamp, err := s.redisConn.Get("timestamp_" + sttString).Uint64()

	if err != nil {
		logs, err = s.db.Query(s.query + " ALLOW FILTERING").Iter().SliceMap()
		if err != nil {
			return nil, err
		}
	} else {
		logs, err = s.db.Query(s.query+" WHERE "+columnName+" > ?  ALLOW FILTERING", timestamp).Iter().SliceMap()
		if err != nil {
			return nil, err
		}
	}

	nameTable, err := getNameTableFromQuery(s.query)
	if err != nil {
		return nil, err
	}

	err = s.db.Query("SELECT MAX(" + columnName + ") FROM " + nameTable + " ALLOW FILTERING").Scan(&timestamp)
	if err != nil {
		return nil, err
	}

	log.Println(stt, timestamp)

	err = s.redisConn.Set("timestamp_"+sttString, timestamp, 10*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	//err = s.redisConn.Set("dns_"+sttString, s.dns, 10*time.Hour).Err()
	//if err != nil {
	//	return nil, err
	//}

	logsJSON, err := json.Marshal(logs)
	if err != nil {
		return nil, err
	}

	return logsJSON, nil
}

func (s *cassandraStory) PushLogBySchedule(writer kafka.IKafkaWriter, ctx context.Context, stt int) {
	logs, err := s.Execute(stt)
	if err != nil {
		log.Fatal(err)
	}
	if len(string(logs)) != 2 {
		log.Println(len(string(logs)))

		err = writer.WriteMessages(ctx, logs, stt)
		if err != nil {
			log.Panic(err)
		}

		log.Println("updated logs on Cassandra", time.Now())
	}
}

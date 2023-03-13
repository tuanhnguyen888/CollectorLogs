package database

import (
	"collector/cache"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"collector/kafka"
)

type IDb interface {
	Raw(sql string, values ...interface{}) IRaw
}

type IRaw interface {
	Rows() (IRow, error)
	Scan(dest interface{}) error
}

type IRow interface {
	Columns() ([]string, error)
	Scan(dest ...any) error
	Close() error
	Next() bool
}

type gormDBImpl struct {
	db *gorm.DB
}

type rawImpl struct {
	raw *gorm.DB
}

type rowsImpl struct {
	row *sql.Rows
}

func (g *gormDBImpl) Raw(sql string, values ...interface{}) IRaw {
	raw := g.db.Raw(sql, values...)
	rawImpl := new(rawImpl)
	rawImpl.raw = raw
	return rawImpl
}

//func(g *gormDBImpl) Row() *sql.Row {
//	return g.db.Row()
//}

func (g *rawImpl) Rows() (IRow, error) {
	rows, err := g.raw.Rows()
	rowsImpl := new(rowsImpl)
	rowsImpl.row = rows
	return rowsImpl, err
}

func (g *rawImpl) Scan(dest interface{}) error {
	err := g.raw.Scan(dest).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *rowsImpl) Columns() ([]string, error) {
	return r.row.Columns()
}
func (r *rowsImpl) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

func (r *rowsImpl) Close() error {
	return r.row.Close()
}

func (r *rowsImpl) Next() bool {
	return r.row.Next()
}

//type IQueryService interface {
//	Query(tx IDb) IDb
//}
//
//type QueryService struct{}
//
//func (s *QueryService) Query(tx IDb) IDb {
//	return tx
//}

type postgresqlStory struct {
	query string
	DB    IDb
	//QueryService IQueryService
	RedisCon cache.IRedisClient
}

func NewGormPgDB(dsn string) IDb {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicln(err)
	}
	return &gormDBImpl{db}
}

func NewPostgres(query string, db IDb, RedisCon cache.IRedisClient) Database {
	return &postgresqlStory{
		query: query,
		DB:    db,
		//QueryService: &QueryService{},
		RedisCon: RedisCon,
	}
}

func (s *postgresqlStory) Execute(stt int) ([]byte, error) {
	var logs []map[string]interface{}
	sttString := strconv.Itoa(stt)
	//dns := s.redisClient.Get("dns_" + sttString).String()
	nameTable, err := getNameTableFromQuery(s.query)
	if err != nil {
		return nil, err
	}

	var columnName string
	rows, err := s.DB.Raw(s.query + " LIMIT 1").Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Lấy thông tin của các trường, kiểu và giá trị từ bản ghi
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}
	if rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal(err)
		}
		for i, col := range columns {
			val := values[i]
			valInt, ok := val.(int64)
			if !ok {
				continue
			}
			if len(strconv.FormatInt(valInt, 10)) == 13 || len(strconv.FormatInt(valInt, 10)) == 11 {
				columnName = col
				break
			}
		}
	} else {
		log.Fatal("No rows found in the result set")
	}

	if columnName == "" {
		return nil, errors.New("Not found column name ")
	}
	//log.Println(columnName)
	// SU dung Redis
	timestamp, err := s.RedisCon.Get("timestamp_" + sttString).Uint64()
	if err != nil {
		err := s.DB.Raw(s.query).Scan(&logs)
		if err != nil {
			return nil, err
		}
	} else {
		err = s.DB.Raw(s.query+" where "+columnName+" > ?", timestamp).Scan(&logs)
		if err != nil {
			return nil, err
		}
	}
	logsJSON, err := json.Marshal(logs)
	if err != nil {
		return nil, err
	}
	err = s.DB.Raw("SELECT MAX(" + columnName + ") FROM " + nameTable).Scan(&timestamp)
	if err != nil {
		return nil, err
	}
	log.Println(stt, timestamp)
	err = s.RedisCon.Set("timestamp_"+sttString, timestamp, 10*time.Hour).Err()
	if err != nil {
		return nil, err
	}
	return logsJSON, nil
}

func (s *postgresqlStory) PushLogBySchedule(writer kafka.IKafkaWriter, ctx context.Context, stt int) {
	logs, err := s.Execute(stt)
	if err != nil {
		log.Fatal(err)
	}
	if len(string(logs)) != 4 {
		log.Println(len(string(logs)))

		err = writer.WriteMessages(ctx, logs, stt)
		if err != nil {
			log.Panic(err)
		}

		log.Println("updated logs on postgres", time.Now())
	}
}

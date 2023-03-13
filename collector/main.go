package main

import (
	"collector/cache"
	"collector/database"
	"collector/kafka"
	"context"
	"fmt"
	"gopkg.in/robfig/cron.v2"
	"log"
	"os"

	"github.com/joho/godotenv"
	"time"
)

type dbConfig struct {
	DbName   string
	Dns      string
	Query    string
	Schedule string
}

func initDatabase(config dbConfig, conn cache.IRedisClient) database.Database {

	switch config.DbName {
	case "postgres":
		IDB := database.NewGormPgDB(config.Dns)
		dataStory := database.NewPostgres(config.Query, IDB, conn)
		return dataStory
	// case "mysql":
	// 	dataStory := database.NewMysql(config.Dns, config.Query, conn)
	// 	return dataStory
	case "mssql":
		IDB := database.NewMssqlPgDB(config.Dns)
		dataStory := database.NewMssql(config.Query, IDB, conn)
		return dataStory
	case "sqlserver":
		IDB := database.NewMssqlPgDB(config.Dns)
		dataStory := database.NewMssql(config.Query, IDB, conn)
		return dataStory
	case "cassandra":
		ICassandra, err := database.NewCassandraConn(config.Dns)
		if err != nil {
			log.Panicln(err)
		}
		dataStory := database.NewCassandra(config.Query, ICassandra, conn)
		return dataStory
	default:
		log.Fatalf("invalid database name")
		return nil
	}
}

func main() {
	redisConn := cache.NewRedisClient()

	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}

	// Get the number of databases
	numDbs := 0
	for {
		_, ok := os.LookupEnv(fmt.Sprintf("DB_NAME_%d", numDbs))
		if !ok {
			break
		}
		numDbs++
	}

	// Create a slice to hold the database configs
	dbConfigs := make([]dbConfig, numDbs)

	// Loop through each database config and add it to the slice
	for i := 0; i < numDbs; i++ {
		name, _ := os.LookupEnv(fmt.Sprintf("DB_NAME_%d", i))
		dsn, _ := os.LookupEnv(fmt.Sprintf("DB_DSN_%d", i))
		query, _ := os.LookupEnv(fmt.Sprintf("DB_QUERY_%d", i))
		schedule, _ := os.LookupEnv(fmt.Sprintf("DB_SCHEDULE_%d", i))
		if err != nil {
			fmt.Printf("Error parsing DB_SCHEDULE_%d: %v\n", i, err)
			os.Exit(1)
		}

		dbConfigs[i] = dbConfig{
			DbName:   name,
			Dns:      dsn,
			Query:    query,
			Schedule: schedule,
		}
	}

	writer := kafka.NewKafkaWriter("logs")
	ctx := context.Background()

	for i := 0; i < numDbs; i++ {
		go func(config dbConfig, index int) {
			data := initDatabase(config, redisConn)

			c := cron.New()
			c.AddFunc(config.Schedule, func() {
				data.PushLogBySchedule(writer, ctx, index)
			})

			c.Start()
		}(dbConfigs[i], i)
	}

	for {
		time.Sleep(time.Hour)
	}

}

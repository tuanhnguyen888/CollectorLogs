package database

import (
	"collector/kafka"
	"context"
	"errors"
	"regexp"
)

type Database interface {
	Execute(stt int) ([]byte, error)
	PushLogBySchedule(writer kafka.IKafkaWriter, ctx context.Context, stt int)
}

func getNameTableFromQuery(query string) (string, error) {
	pattern := regexp.MustCompile("FROM\\s+(\\S+)")
	match := pattern.FindStringSubmatch(query)
	if len(match) > 1 {
		return match[1], nil
	} else {
		return "", errors.New("Not found name table database ")
	}

}

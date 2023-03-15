package kafka_test

import (
	"context"
	"testing"

	"collector/kafka"
	"github.com/stretchr/testify/assert"
)

func TestNewKafkaWriter(t *testing.T) {
	topic := "test"

	kafkaWriter := kafka.NewKafkaWriter(topic)

	assert.IsType(t, &kafka.Writer{}, kafkaWriter)
	assert.Len(t, kafkaWriter.(*kafka.Writer).Writer, 1)
	assert.Equal(t, topic, kafkaWriter.(*kafka.Writer).Writer[0].Topic)
}

func TestWriteMessages(t *testing.T) {
	// create a new kafka writer instance
	kafkaWriter := kafka.NewKafkaWriter("test")

	// set up some test data
	testData := []byte("test data")

	// create a new context
	ctx := context.Background()

	// call WriteMessages method to write test data to kafka
	err := kafkaWriter.WriteMessages(ctx, testData, 0)
	assert.NoError(t, err, "unexpected error")
}

package kafka

import (
	"context"
	"log"

	kafkago "github.com/segmentio/kafka-go"
)

type IKafkaWriter interface {
	WriteMessages(ctx context.Context, logEvent []byte, partitionNumber int) error
}

type Writer struct {
	Writer []*kafkago.Writer
}

func NewKafkaWriter(topic string) IKafkaWriter {
	brokers := []string{"127.0.0.1:9092"}
	writers := make([]*kafkago.Writer, len(brokers))
	for i, broker := range brokers {
		writers[i] = &kafkago.Writer{
			Addr:                   kafkago.TCP(broker),
			Topic:                  topic,
			AllowAutoTopicCreation: true,
			Balancer:               &kafkago.LeastBytes{},
		}
	}

	return &Writer{
		Writer: writers,
	}
}

func (k *Writer) WriteMessages(ctx context.Context, logEvent []byte, partitionNumber int) error {

	messages := kafkago.Message{}

	messages.Value = logEvent

	for _, writer := range k.Writer {
		err := writer.WriteMessages(ctx,
			kafkago.Message{
				Key:       nil,
				Value:     messages.Value,
				Partition: partitionNumber,
			},
		)
		if err != nil {
			log.Fatalf("Error writing message: %v", err)
		}
	}

	return nil
}

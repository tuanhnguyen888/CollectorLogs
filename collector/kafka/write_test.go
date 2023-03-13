// writer_test.go

package kafka

import (
	"collector/mocks"
	"context"
	kafkago "github.com/segmentio/kafka-go"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWriter_WriteMessages(t *testing.T) {
	// Tạo mock
	mockWriter := mocks.NewIKafkaWriter(t)

	// Thiết lập expectation cho mock
	logEvent := []byte("log message")
	partitionNumber := 0
	ctx := context.Background()
	mockWriter.On("WriteMessages", mock.Anything, logEvent, partitionNumber).Return(nil)

	// Tạo writer
	writer := &Writer{
		Writer: []*kafkago.Writer{mockWriter},
	}

	// Gọi hàm WriteMessages
	err := writer.WriteMessages(ctx, logEvent, partitionNumber)

	// Kiểm tra kết quả
	assert.NoError(t, err)

	// Kiểm tra expectation đã được gọi
	mockWriter.AssertExpectations(t)
}

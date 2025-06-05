//go:build e2e

package mqx

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MockUserEvent struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestNewGeneralProducer(t *testing.T) {
	t.Skip()
	suite.Run(t, new(TestGeneralProducerTestSuite))
}

type TestGeneralProducerTestSuite struct {
	suite.Suite
}

func (s *TestGeneralProducerTestSuite) TestProduceAndConsume() {
	t := s.T()
	addr := "localhost:9092"
	topic := "mock_user_events"

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": addr,
	})
	assert.NoError(t, err)

	producer, err := NewGeneralProducer[MockUserEvent](kafkaProducer, topic)
	assert.NoError(t, err)

	defer producer.Close()

	expected := MockUserEvent{
		Name: "alex",
		Age:  18,
	}
	err = producer.Produce(t.Context(), expected)

	assert.NoError(t, err)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  addr,
		"group.id":           fmt.Sprintf("test-%d", time.Now().UnixNano()),
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "false",
	})
	assert.NoError(t, err)

	err = consumer.SubscribeTopics([]string{topic}, nil)
	assert.NoError(t, err)

	message, err := consumer.ReadMessage(10 * time.Second)
	assert.NoError(t, err)

	var actual MockUserEvent
	err = json.Unmarshal(message.Value, &actual)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)

	// 没提交，再次消费相同消息
	err = json.Unmarshal(message.Value, &actual)
	assert.NoError(t, err)

	ps, err := consumer.CommitMessage(message)
	assert.NoError(t, err)

	const i = 0
	assert.Equal(t, topic, *ps[i].Topic)
	assert.NoError(t, ps[i].Error)

	// 提交后，再次消费，应该超时
	_, err = consumer.ReadMessage(10 * time.Second)
	assert.Error(t, err)
}

package confluentinc_llm

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/config"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/logger"
	"strings"
	"time"
)

const (
	TimeOut = 500
)

type LlmConsumer struct {
	c   *kafka.Consumer
	log logger.Logger
}

func initConsumer(conf *config.KafkaConfig, groupId string) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(conf.Addr, ","),
		"auto.offset.reset":  "latest",
		"group.id":           groupId,
		"enable.auto.commit": true,
	})

	if err != nil {
		panic(err)
	}

	return c
}

func NewLlmConsumer(conf *config.KafkaConfig, groupId string) *LlmConsumer {
	c := initConsumer(conf, groupId)
	return &LlmConsumer{c: c, log: logger.NewDefaultLogger()}
}

func (lc *LlmConsumer) Subscribe(topics []string) {
	err := lc.c.SubscribeTopics(topics, nil)
	if err != nil {
		panic(err)
	}
}

func (lc *LlmConsumer) Consume() (string, []byte) {
	m, err := lc.c.ReadMessage(TimeOut * time.Millisecond)
	if err != nil {
		var kErr kafka.Error
		if errors.As(err, &kErr) && kErr.Code() == kafka.ErrTimedOut {
			return "", nil
		}

		lc.log.Error("consume error", logger.Error(err))
		return "", nil
	}

	return *m.TopicPartition.Topic, m.Value
}

func (lc *LlmConsumer) WrapLogger(log logger.Logger) {
	lc.log = log
}

func (lc *LlmConsumer) Close() {
	if err := lc.c.Close(); err != nil {
		lc.log.Error("close consumer error", logger.Error(err))
	}
}

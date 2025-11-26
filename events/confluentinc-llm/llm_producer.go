package confluentinc_llm

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/logger"
)

const KafkaFailed = "LLMAuditProduceFailed"

type FailedMsg struct {
	Topic string `json:"topic"`
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

type LlmProducer struct {
	p      *kafka.Producer
	failed chan *kafka.Message
	db     *redis.Client
	log    logger.Logger
	wg     sync.WaitGroup
}

func NewLlmProducer(p *kafka.Producer, db *redis.Client) *LlmProducer {
	lp := &LlmProducer{
		p:      p,
		db:     db,
		log:    logger.NewDefaultLogger(),
		failed: make(chan *kafka.Message),
	}

	go func() {
		for e := range lp.p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					lp.log.Error("发送失败，加入重试队列:",
						logger.Error(ev.TopicPartition.Error))

					copyMsg := *ev
					lp.failed <- &copyMsg
				}
			}
		}
	}()

	lp.wg.Add(1)
	go func() {
		defer lp.wg.Done()

		for msg := range lp.failed {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

			b, _ := json.Marshal(FailedMsg{
				Topic: *msg.TopicPartition.Topic,
				Key:   msg.Key,
				Value: msg.Value,
			})

			msgID := fmt.Sprintf("%s:%d", *msg.TopicPartition.Topic, msg.Key)

			if lp.db.SetNX(context.Background(), msgID, 1, 24*time.Hour).Val() {
				lp.db.RPush(ctx, KafkaFailed, b)
			}

			cancel()
		}
	}()
	return lp
}

func (lp *LlmProducer) WrapLogger(log logger.Logger) {
	lp.log = log
}

func (lp *LlmProducer) Produce(topic string, key, data []byte) {
	err := lp.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          data,
	}, nil)

	if err != nil {
		lp.log.Error("Produce immediate error", logger.Error(err))

		msgID := fmt.Sprintf("%s:%d", topic, key)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if lp.db.SetNX(ctx, msgID, 1, 24*time.Hour).Val() {
			b, _ := json.Marshal(FailedMsg{
				Topic: topic,
				Key:   key,
				Value: data,
			})
			lp.db.RPush(ctx, KafkaFailed, b)
		}
	}
}

func (lp *LlmProducer) Close(timeout time.Duration) {
	lp.p.Flush(int(timeout.Milliseconds()))
	lp.p.Close()

	close(lp.failed)

	lp.wg.Wait()
}

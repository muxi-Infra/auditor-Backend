package ioc

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/config"
)

func InitProducer(conf *config.KafkaConfig) *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     conf.User,
		"sasl.password":     conf.Password,
		"bootstrap.servers": strings.Join(conf.Addr, ","), // 这里连宿主机映射的端口
		"retries":           5,                            // 内部重试次数
		"retry.backoff.ms":  500,                          // 重试间隔
		"acks":              "all",                        // 确保 leader 和 ISR 都收到才算成功
	})

	if err != nil {
		panic(err)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Println("发送失败:", ev.TopicPartition.Error)
				}
			}
		}
	}()

	return p
}

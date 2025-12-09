package events

import (
	"time"

	"github.com/muxi-Infra/auditor-Backend/pkg/logger"
)

type Producer interface {
	Produce(topic string, key, data []byte)
	WrapLogger(log logger.Logger)
	Close(timeout time.Duration)
}

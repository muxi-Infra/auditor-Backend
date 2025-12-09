package events

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/logger"
	"time"
)

type Producer interface {
	Produce(topic string, key, data []byte)
	WrapLogger(log logger.Logger)
	Close(timeout time.Duration)
}

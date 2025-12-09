package events

import "github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/logger"

type Consumer interface {
	Subscribe(topics []string)
	Consume() (string, []byte)
	Close()
	WrapLogger(log logger.Logger)
}

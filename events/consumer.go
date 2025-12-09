package events

import "github.com/muxi-Infra/auditor-Backend/pkg/logger"

type Consumer interface {
	Subscribe(topics []string)
	Consume() (string, []byte)
	Close()
	WrapLogger(log logger.Logger)
}

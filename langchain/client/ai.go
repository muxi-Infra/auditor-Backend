package client

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/client/core/ollamas"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/config"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/model"
)

type AuditAIClient interface {
	SendMessage(prompt string) (model.AuditResult, error)
}

func Connect(conf *config.MuxiAI) AuditAIClient {
	switch conf.Type {
	case config.Ollama:
		client := ollamas.NewOllamaClient(conf.URL, conf.Model)
		return client
	case config.OpenAI:
		return nil
	default:
		panic("illegal AI config")
	}
}

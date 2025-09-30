package client

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/config"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/model"
)

type MuxiAIClient interface {
	SendMessage(prompt string) (model.AIResponse, error)
}

func Connect(conf *config.MuxiAI) {
	switch conf.Type {
	case config.Ollama:

	case config.OpenAI:
	}
}

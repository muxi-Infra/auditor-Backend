package client

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/client/core/ali"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/client/core/ollamas"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/config"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/model"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/logger"
)

type AuditAIClient interface {
	SendMessage(content string, pics []string) (model.AuditResult, error)
	WrapLogger(logger logger.Logger)
	Transform(role string, contents response.Contents) (string, []string)
}

func AuditAIConnect(conf *config.MuxiAI) AuditAIClient {
	switch conf.Type {
	case config.Ollama:
		client := ollamas.NewOllamaClient(conf.URL, conf.Model)
		return client
	case config.OpenAI:
		return nil
	case config.Alibaba:
		return ali.NewAlClient(conf.AccessKeyID, conf.AccessKeySecret, conf.Region, conf.Endpoint)
	default:
		panic("illegal AI config")
	}
}

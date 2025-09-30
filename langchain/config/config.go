package config

import "github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/viperx"

type MuxiAIType string

const (
	Ollama MuxiAIType = "ollama"
	OpenAI MuxiAIType = "openai"
)

type MuxiAI struct {
	Type   MuxiAIType `yaml:"type"`
	URL    string     `yaml:"url"`
	Model  string     `yaml:"model"`
	ApiKey string     `yaml:"api_key"`
}

func NewMuxiAIConf(s *viperx.VipperSetting) *MuxiAI {
	var aiConf = &MuxiAI{}
	err := s.ReadSection("AI", &aiConf)
	if err != nil {
		panic(err)
	}
	return aiConf
}
